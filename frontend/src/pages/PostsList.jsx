import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { apiUrl } from '../api'

const PAGE_SIZE = 10

function PostsList() {
  const [posts, setPosts] = useState([])
  const [loading, setLoading] = useState(true)
  const [loadingMore, setLoadingMore] = useState(false)
  const [error, setError] = useState(null)
  const [hasMore, setHasMore] = useState(true)
  const [offset, setOffset] = useState(0)

  useEffect(() => {
    document.title = 'Posts - Blog'
  }, [])

  useEffect(() => {
    setLoading(true)
    setError(null)
    const url = apiUrl(`/api/posts?limit=${PAGE_SIZE}&offset=0`)
    const controller = new AbortController()
    const timeoutId = setTimeout(() => controller.abort(), 10000)
    fetch(url, { signal: controller.signal })
      .then((res) => {
        clearTimeout(timeoutId)
        if (!res.ok) {
          return res.json().then((data) => {
            throw new Error(data?.error || `HTTP ${res.status}`)
          }).catch((e) => {
            if (e instanceof Error && e.message.startsWith('HTTP')) throw e
            throw new Error(res.statusText || `HTTP ${res.status}`)
          })
        }
        const ct = res.headers.get('content-type')
        if (!ct || !ct.includes('application/json')) {
          throw new Error(`Backend returned non-JSON (Content-Type: ${ct || 'none'}), check API URL: ${url}`)
        }
        return res.json()
      })
      .then((data) => {
        const list = data.posts || []
        setPosts(list)
        setOffset(list.length)
        setHasMore(list.length >= PAGE_SIZE)
        setLoading(false)
      })
      .catch((err) => {
        clearTimeout(timeoutId)
        let msg = err.message || 'Network error'
        if (err.name === 'AbortError') {
          msg = `Request timeout. Ensure backend is running: ${apiUrl('/health')}`
        } else if (typeof msg === 'string' && (msg.includes('Failed to fetch') || msg.includes('Connection refused') || msg.includes('NetworkError'))) {
          msg = `Connection refused — backend not running? Start with: ./scripts/start-dev.sh (or run backend in another terminal: cd backend && go run main.go)`
        }
        setError(msg)
        setLoading(false)
      })
  }, [])

  useEffect(() => {
    // SSE: this request stays open (shows as "pending" in DevTools) — that's expected
    const url = apiUrl('/api/events')
    const es = new EventSource(url)
    es.addEventListener('sync_completed', () => {
      fetch(apiUrl(`/api/posts?limit=${PAGE_SIZE}&offset=0`))
        .then((res) => res.ok ? res.json() : Promise.reject(new Error('refetch failed')))
        .then((data) => {
          const list = data.posts || []
          setPosts(list)
          setOffset(list.length)
          setHasMore(list.length >= PAGE_SIZE)
        })
        .catch(() => {})
    })
    return () => es.close()
  }, [])

  function loadMore() {
    if (loadingMore || !hasMore) return
    setLoadingMore(true)
    fetch(apiUrl(`/api/posts?limit=${PAGE_SIZE}&offset=${offset}`))
      .then((res) => res.json())
      .then((data) => {
        const list = data.posts || []
        setPosts((prev) => [...prev, ...list])
        setOffset((prev) => prev + list.length)
        setHasMore(list.length >= PAGE_SIZE)
        setLoadingMore(false)
      })
      .catch(() => setLoadingMore(false))
  }

  if (loading) {
    return (
      <div className="text-center py-12">
        <div className="text-gray-500">Loading...</div>
        <p className="text-sm text-gray-400 mt-2">If no response, ensure backend is running: <a href={apiUrl('/health')} target="_blank" rel="noopener noreferrer" className="text-blue-600">{apiUrl('/health')}</a></p>
      </div>
    )
  }

  if (error) {
    return (
      <div className="text-center py-12">
        <div className="text-red-500">Load failed: {error}</div>
      </div>
    )
  }

  if (posts.length === 0) {
    return (
      <div className="text-center py-12">
        <p className="text-gray-500">No posts yet.</p>
      </div>
    )
  }

  return (
    <div>
      <h1 className="text-3xl font-bold mb-8 text-gray-900">Posts</h1>
      <div className="space-y-6">
        {posts.map((post) => (
          <article
            key={post.id}
            className="bg-white rounded-lg shadow-sm p-6 hover:shadow-md transition-shadow"
          >
            <Link to={`/posts/${post.slug}`}>
              <h2 className="text-2xl font-semibold text-gray-900 mb-2 hover:text-blue-600 transition-colors">
                {post.title || 'Untitled'}
              </h2>
            </Link>
            {post.summary && (
              <p className="text-gray-600 mb-4 line-clamp-2">{post.summary}</p>
            )}
            <div className="flex items-center justify-between text-sm text-gray-500">
              <div className="flex items-center space-x-4">
                <time dateTime={post.published_at}>
                  {new Date(post.published_at).toLocaleDateString('en-US', {
                    year: 'numeric',
                    month: 'long',
                    day: 'numeric',
                  })}
                </time>
                {post.category && (
                  <span className="px-2 py-1 bg-gray-100 rounded text-gray-700">
                    {post.category}
                  </span>
                )}
              </div>
              <Link
                to={`/posts/${post.slug}`}
                className="text-blue-600 hover:text-blue-800 font-medium"
              >
                Read more →
              </Link>
            </div>
          </article>
        ))}
      </div>
      {hasMore && posts.length > 0 && (
        <div className="mt-8 text-center">
          <button
            onClick={loadMore}
            disabled={loadingMore}
            className="px-6 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50 disabled:opacity-50 transition-colors"
          >
            {loadingMore ? 'Loading...' : 'Load more'}
          </button>
        </div>
      )}
    </div>
  )
}

export default PostsList
