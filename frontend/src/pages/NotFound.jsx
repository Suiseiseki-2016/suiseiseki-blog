import { useEffect } from 'react'
import { Link } from 'react-router-dom'

function NotFound() {
  useEffect(() => {
    document.title = '404 - 页面不存在 - Blog'
    return () => { document.title = 'Blog' }
  }, [])

  return (
    <div className="text-center py-16">
      <h1 className="text-6xl font-bold text-gray-300 mb-4">404</h1>
      <p className="text-xl text-gray-600 mb-8">页面不存在</p>
      <Link
        to="/"
        className="inline-block px-6 py-3 bg-gray-900 text-white rounded-lg hover:bg-gray-800 transition-colors"
      >
        返回首页
      </Link>
    </div>
  )
}

export default NotFound
