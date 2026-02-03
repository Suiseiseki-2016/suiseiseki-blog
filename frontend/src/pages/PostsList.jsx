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
    document.title = '文章列表 - Blog'
  }, [])

  // 初始加载
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
          throw new Error(`后端返回非 JSON (Content-Type: ${ct || '无'})，请确认 API 地址正确: ${url}`)
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
        const msg = err.name === 'AbortError'
          ? '请求超时，请确认后端已启动: http://localhost:8080/health'
          : (err.message || '网络错误')
        setError(msg)
        setLoading(false)
      })
  }, [])

  // 监听后端同步完成事件，静默刷新列表（不整页重载）
  useEffect(() => {
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
        <div className="text-gray-500">加载中...</div>
        <p className="text-sm text-gray-400 mt-2">若长时间无响应，请确认后端已启动：<a href="http://localhost:8080/health" target="_blank" rel="noopener noreferrer" className="text-blue-600">http://localhost:8080/health</a></p>
      </div>
    )
  }

  if (error) {
    return (
      <div className="text-center py-12">
        <div className="text-red-500">加载失败: {error}</div>
      </div>
    )
  }

  if (posts.length === 0) {
    return (
      <div className="text-center py-12">
        <p className="text-gray-500">暂无文章</p>
      </div>
    )
  }

  return (
    <div>
      <h1 className="text-3xl font-bold mb-8 text-gray-900">文章列表</h1>
      <div className="space-y-6">
        {posts.map((post) => (
          <article
            key={post.id}
            className="bg-white rounded-lg shadow-sm p-6 hover:shadow-md transition-shadow"
          >
            <Link to={`/posts/${post.slug}`}>
              <h2 className="text-2xl font-semibold text-gray-900 mb-2 hover:text-blue-600 transition-colors">
                {post.title || '无标题'}
              </h2>
            </Link>
            {post.summary && (
              <p className="text-gray-600 mb-4 line-clamp-2">{post.summary}</p>
            )}
            <div className="flex items-center justify-between text-sm text-gray-500">
              <div className="flex items-center space-x-4">
                <time dateTime={post.published_at}>
                  {new Date(post.published_at).toLocaleDateString('zh-CN', {
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
                阅读更多 →
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
            {loadingMore ? '加载中...' : '加载更多'}
          </button>
        </div>
      )}
    </div>
  )
}

export default PostsList
