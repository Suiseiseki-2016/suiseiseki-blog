import { useEffect, useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import { apiUrl } from '../api'

function PostDetail() {
  const { slug } = useParams()
  const [post, setPost] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)

  useEffect(() => {
    if (post?.title) document.title = `${post.title} - Blog`
    return () => { document.title = 'Blog' }
  }, [post])

  useEffect(() => {
    fetch(apiUrl(`/api/posts/${slug}`))
      .then((res) => {
        if (!res.ok) {
          throw new Error('文章不存在')
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

  if (loading) {
    return (
      <div className="text-center py-12">
        <div className="text-gray-500">加载中...</div>
      </div>
    )
  }

  if (error || !post) {
    return (
      <div className="text-center py-12">
        <div className="text-red-500 mb-4">加载失败: {error || '文章不存在'}</div>
        <Link to="/" className="text-blue-600 hover:text-blue-800">
          返回首页
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
          ← 返回列表
        </Link>
        <h1 className="text-4xl font-bold text-gray-900 mb-4">{post.title}</h1>
        <div className="flex items-center space-x-4 text-sm text-gray-600">
          <time dateTime={post.published_at}>
            {new Date(post.published_at).toLocaleDateString('zh-CN', {
              year: 'numeric',
              month: 'long',
              day: 'numeric',
            })}
          </time>
          {post.category && (
            <span className="px-3 py-1 bg-gray-100 rounded text-gray-700">
              {post.category}
            </span>
          )}
        </div>
      </header>
      <div
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
