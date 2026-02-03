import { Link } from 'react-router-dom'

function Layout({ children }) {
  return (
    <div className="min-h-screen bg-gray-50">
      <nav className="bg-white shadow-sm">
        <div className="max-w-4xl mx-auto px-4 py-4">
          <div className="flex justify-between items-center">
            <Link to="/" className="text-xl font-bold text-gray-900 hover:text-gray-700">
              Blog
            </Link>
            <div className="space-x-4">
              <Link
                to="/"
                className="text-gray-600 hover:text-gray-900 transition-colors"
              >
                Posts
              </Link>
              <Link
                to="/resume"
                className="text-gray-600 hover:text-gray-900 transition-colors"
              >
                Resume
              </Link>
            </div>
          </div>
        </div>
      </nav>
      <main className="max-w-4xl mx-auto px-4 py-8">
        {children}
      </main>
      <footer className="bg-white border-t mt-12">
        <div className="max-w-4xl mx-auto px-4 py-6 text-center text-gray-600 text-sm">
          © {new Date().getFullYear()} Blog. Built with Go + React.
        </div>
      </footer>
    </div>
  )
}

export default Layout
