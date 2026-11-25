import { useEffect, useMemo, useRef, useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import 'katex/dist/katex.min.css'
import { apiUrl } from '../api'

function stripLeadingDuplicateH1Html(html, title) {
  if (!html || !title || typeof DOMParser === 'undefined') return html
  const doc = new DOMParser().parseFromString(html, 'text/html')
  const first = doc.body.firstElementChild
  if (first && first.tagName === 'H1') {
    const norm = (s) => (s || '').replace(/\s+/g, ' ').trim().toLowerCase()
    if (norm(first.textContent) === norm(title)) {
      first.remove()
      return doc.body.innerHTML
    }
  }
  return html
}

const DISPLAY_MATH_RE = /\$\$([\s\S]+?)\$\$/g
const INLINE_MATH_RE = /(?<!\\|\$)\$(?!\s|\$)([^$\n\u4e00-\u9fff]+?)(?<!\s|\\)\$(?!\$|\d)/g

function renderInlineMathSegment(text, katex) {
  const fragment = document.createDocumentFragment()
  let lastIdx = 0
  INLINE_MATH_RE.lastIndex = 0
  let match
  while ((match = INLINE_MATH_RE.exec(text)) !== null) {
    const expr = match[1].trim()
    const isLikelyCurrency = /^\d/.test(expr) && !/[=+\-*/^_\\<>]/.test(expr)
    if (!expr || /[\u4e00-\u9fff]/.test(expr) || isLikelyCurrency) {
      continue
    }
    if (match.index > lastIdx) {
      fragment.appendChild(document.createTextNode(text.slice(lastIdx, match.index)))
    }
    const span = document.createElement('span')
    span.innerHTML = katex.renderToString(expr, { displayMode: false, throwOnError: false })
    fragment.appendChild(span)
    lastIdx = match.index + match[0].length
  }
  if (lastIdx === 0) return null
  if (lastIdx < text.length) {
    fragment.appendChild(document.createTextNode(text.slice(lastIdx)))
  }
  return fragment
}

function renderMathInElement(root, katex) {
  const walker = document.createTreeWalker(root, NodeFilter.SHOW_TEXT)
  const textNodes = []
  let node
  while ((node = walker.nextNode())) {
    if (!node.nodeValue || !node.nodeValue.includes('$')) continue
    if (node.parentElement?.closest('pre, code, script, style, textarea, .mermaid, .katex, a')) continue
    textNodes.push(node)
  }

  for (const textNode of textNodes) {
    const raw = textNode.nodeValue
    const fragment = document.createDocumentFragment()
    let changed = false
    let lastIdx = 0
    DISPLAY_MATH_RE.lastIndex = 0
    let match

    while ((match = DISPLAY_MATH_RE.exec(raw)) !== null) {
      changed = true
      if (match.index > lastIdx) {
        const prefix = raw.slice(lastIdx, match.index)
        const inlineFrag = renderInlineMathSegment(prefix, katex)
        fragment.appendChild(inlineFrag || document.createTextNode(prefix))
      }
      const div = document.createElement('div')
      div.innerHTML = katex.renderToString(match[1].trim(), { displayMode: true, throwOnError: false })
      fragment.appendChild(div)
      lastIdx = match.index + match[0].length
    }

    const tail = raw.slice(lastIdx)
    if (tail) {
      const inlineFrag = renderInlineMathSegment(tail, katex)
      if (inlineFrag) {
        changed = true
        fragment.appendChild(inlineFrag)
      } else if (changed) {
        fragment.appendChild(document.createTextNode(tail))
      }
    }

    if (changed && textNode.parentNode) {
      textNode.parentNode.replaceChild(fragment, textNode)
    }
  }
}

function PostDetail() {
  const { slug } = useParams()
  const [post, setPost] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const contentRef = useRef(null)

  const cleanContent = useMemo(
    () => stripLeadingDuplicateH1Html(post?.content || '', post?.title || ''),
    [post?.content, post?.title]
  )

  useEffect(() => {
    if (post?.title) document.title = `${post.title} - Di Wu`
    return () => { document.title = 'Di Wu' }
  }, [post])

  useEffect(() => {
    setLoading(true)
    setError(null)
    fetch(apiUrl(`/api/posts/${slug}`))
      .then((res) => {
        if (!res.ok) {
            throw new Error('Post not found')
        }
        return res.json()
      })
      .then((data) => {
        setPost(data)
        setLoading(false)
      })
      .catch((err) => {
        setError(err.message)
        setLoading(false)
      })
  }, [slug])

  useEffect(() => {
    if (!cleanContent || !contentRef.current) return
    const mermaidEls = [...contentRef.current.querySelectorAll('.mermaid')].filter(
      (el) => !el.querySelector('svg')
    )
    if (!mermaidEls.length) return
    import('mermaid').then(({ default: mermaid }) => {
      mermaid.initialize({ startOnLoad: false, theme: 'default' })
      mermaid.run({ nodes: mermaidEls })
    }).catch(console.error)
  }, [cleanContent])

  useEffect(() => {
    if (!cleanContent || !contentRef.current) return
    if (!contentRef.current.textContent?.includes('$')) return
    let cancelled = false
    import('katex').then(({ default: katex }) => {
      if (cancelled || !contentRef.current) return
      renderMathInElement(contentRef.current, katex)
    }).catch(console.error)
    return () => { cancelled = true }
  }, [cleanContent])

  if (loading) {
    return (
      <div className="text-center py-12">
        <div className="text-gray-500">Loading...</div>
      </div>
    )
  }

  if (error || !post) {
    return (
      <div className="text-center py-12">
        <div className="text-red-500 mb-4">Load failed: {error || 'Post not found'}</div>
        <Link to="/" className="text-blue-600 hover:text-blue-800">
          Back to home
        </Link>
      </div>
    )
  }

  return (
    <article className="bg-white rounded-lg shadow-sm p-5 sm:p-8">
      <header className="mb-8 pb-6 border-b">
        <Link
          to="/"
          className="text-blue-600 hover:text-blue-800 mb-4 inline-block text-sm"
        >
          ← Back to list
        </Link>
        <h1 className="text-2xl sm:text-4xl font-bold text-gray-900 mb-4">{post.title}</h1>
        <div className="flex flex-wrap items-center gap-2 text-sm text-gray-600">
          <time dateTime={post.published_at} className="whitespace-nowrap shrink-0 mr-1">
            {(() => {
              if (!post.published_at) return ''
              const dateStr = post.published_at.split('T')[0]
              const [y, m, d] = dateStr.split('-').map(Number)
              return new Date(y, m - 1, d).toLocaleDateString('zh-CN', { year: 'numeric', month: 'long', day: 'numeric' })
            })()}
          </time>
          {post.tags && post.tags.map((tag) => (
            <span key={tag} className="px-3 py-1 bg-gray-100 rounded text-gray-700 whitespace-nowrap shrink-0">
              {tag}
            </span>
          ))}
        </div>
      </header>
      <div
        ref={contentRef}
        className="prose prose-lg max-w-none"
        dangerouslySetInnerHTML={{ __html: cleanContent }}
        style={{
          lineHeight: '1.8',
        }}
      />
    </article>
  )
}

export default PostDetail
