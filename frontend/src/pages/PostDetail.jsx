import { useEffect, useRef, useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import { apiUrl } from '../api'

function PostDetail() {
  const { slug } = useParams()
  const [post, setPost] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const contentRef = useRef(null)

  useEffect(() => {
    if (post?.title) document.title = `${post.title} - Blog`
    return () => { document.title = 'Blog' }
  }, [post])

  useEffect(() => {
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
    if (!post?.content || !contentRef.current) return
    const mermaidEls = contentRef.current.querySelectorAll('.mermaid')
    if (!mermaidEls.length) return
    import('mermaid').then(({ default: mermaid }) => {
      mermaid.initialize({ startOnLoad: false, theme: 'default' })
      mermaid.run({ nodes: mermaidEls })
    }).catch(console.error)
  }, [post])

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
    <article className="bg-white rounded-lg shadow-sm p-8">
      <header className="mb-8 pb-6 border-b">
        <Link
          to="/"
          className="text-blue-600 hover:text-blue-800 mb-4 inline-block text-sm"
        >
          ← Back to list
        </Link>
        <h1 className="text-4xl font-bold text-gray-900 mb-4">{post.title}</h1>
        <div className="flex items-center space-x-4 text-sm text-gray-600">
          <time dateTime={post.published_at}>
            {(() => {
              const [y, m, d] = post.published_at.split('T')[0].split('-').map(Number)
              return new Date(y, m - 1, d).toLocaleDateString('zh-CN', { year: 'numeric', month: 'long', day: 'numeric' })
            })()}
          </time>
          {post.tags && post.tags.map((tag) => (
            <span key={tag} className="px-3 py-1 bg-gray-100 rounded text-gray-700">
              {tag}
            </span>
          ))}
        </div>
      </header>
      <div
        ref={contentRef}
        className="prose prose-lg max-w-none"
        dangerouslySetInnerHTML={{ __html: post.content }}
        style={{
          lineHeight: '1.8',
        }}
      />
    </article>
  )
}

export default PostDetail
