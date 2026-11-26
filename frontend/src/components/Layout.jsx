import { Link, useLocation } from 'react-router-dom'

function Layout({ children }) {
  const { pathname } = useLocation()
  const isPosts = pathname === '/' || pathname.startsWith('/posts')
  const isResume = pathname === '/resume'

  return (
    <div className="min-h-screen bg-gray-50 flex flex-col">
      <nav className="bg-white shadow-sm">
        <div className="max-w-4xl mx-auto px-4 py-4">
          <div className="flex justify-between items-center gap-4">
            <Link to="/" className="text-xl font-bold text-gray-900 hover:text-gray-700 shrink-0">
              Di Wu
            </Link>
            <div className="flex items-center flex-wrap justify-end gap-3 sm:gap-5 text-sm sm:text-base">
              <Link
                to="/"
                className={`transition-colors ${isPosts ? 'text-gray-900 font-semibold' : 'text-gray-600 hover:text-gray-900'}`}
              >
                Posts
              </Link>
              <Link
                to="/resume"
                className={`transition-colors ${isResume ? 'text-gray-900 font-semibold' : 'text-gray-600 hover:text-gray-900'}`}
              >
                Resume
              </Link>
              <a
                href="https://github.com/Suiseiseki-2016"
                target="_blank"
                rel="noopener noreferrer"
                className="text-gray-600 hover:text-gray-900 transition-colors"
              >
                GitHub
              </a>
              <a
                href="https://pku.aeoluswu.info/"
                target="_blank"
                rel="noopener noreferrer"
                title="PKU All in One — 北京大学（Peking University）教学网学习工作台"
                className="text-gray-600 hover:text-gray-900 transition-colors whitespace-nowrap"
              >
                PKU All in One
              </a>
            </div>
          </div>
        </div>
      </nav>
      <main className="max-w-4xl w-full mx-auto px-4 py-8 flex-1">
        {children}
      </main>
      <footer className="bg-white border-t mt-12">
        <div className="max-w-4xl mx-auto px-4 py-6 text-center text-gray-500 text-sm">
          © {new Date().getFullYear()} 吴迪 (Di Wu) · 探索 AI Agent 与系统安全的边界 ·{' '}
          <a
            href="https://pku.aeoluswu.info/"
            target="_blank"
            rel="noopener noreferrer"
            className="hover:text-gray-900 underline underline-offset-2"
          >
            PKU All in One（北京大学教学网学习工作台）
          </a>
        </div>
      </footer>
    </div>
  )
}

export default Layout
