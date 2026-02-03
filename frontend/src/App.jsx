import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'
import Layout from './components/Layout'
import PostsList from './pages/PostsList'
import PostDetail from './pages/PostDetail'
import Resume from './pages/Resume'
import NotFound from './pages/NotFound'

function App() {
  return (
    <Router>
      <Layout>
        <Routes>
          <Route path="/" element={<PostsList />} />
          <Route path="/posts/:slug" element={<PostDetail />} />
          <Route path="/resume" element={<Resume />} />
          <Route path="*" element={<NotFound />} />
        </Routes>
      </Layout>
    </Router>
  )
}

export default App
