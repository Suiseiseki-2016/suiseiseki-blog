import { useEffect } from 'react'

function Resume() {
  useEffect(() => {
    document.title = 'Resume - Blog'
    return () => { document.title = 'Blog' }
  }, [])

  return (
    <div className="bg-white rounded-lg shadow-sm p-8 print:p-4">
      {/* Print styles */}
      <style>
        {`
          @media print {
            body {
              background: white;
            }
            nav, footer {
              display: none;
            }
            .print-hidden {
              display: none;
            }
            .resume-container {
              max-width: 100%;
              padding: 0;
            }
            .page-break {
              page-break-after: always;
            }
          }
        `}
      </style>

      <div className="resume-container max-w-3xl mx-auto">
        <header className="mb-8 text-center print:text-left">
          <h1 className="text-4xl font-bold text-gray-900 mb-2">Your Name</h1>
          <div className="text-gray-600 space-x-4">
            <span>Email: your.email@example.com</span>
            <span>|</span>
            <span>Phone: +86 138-0000-0000</span>
          </div>
          <div className="text-gray-600 mt-2">
            <span>GitHub: github.com/yourusername</span>
            <span className="mx-2">|</span>
            <span>Website: yourwebsite.com</span>
          </div>
        </header>

        <section className="mb-8">
          <h2 className="text-2xl font-bold text-gray-900 mb-4 border-b pb-2">
            Summary
          </h2>
          <p className="text-gray-700 leading-relaxed">
            Your personal summary: background, skills, career goals.
          </p>
        </section>

        <section className="mb-8">
          <h2 className="text-2xl font-bold text-gray-900 mb-4 border-b pb-2">
            Experience
          </h2>
          <div className="space-y-6">
            <div>
              <div className="flex justify-between items-start mb-2">
                <h3 className="text-xl font-semibold text-gray-900">
                  Job Title
                </h3>
                <span className="text-gray-600">2020.01 - Present</span>
              </div>
              <div className="text-gray-700 font-medium mb-1">Company Name</div>
              <ul className="text-gray-700 list-disc list-inside space-y-1">
                <li>Description 1</li>
                <li>Description 2</li>
                <li>Description 3</li>
              </ul>
            </div>
          </div>
        </section>

        <section className="mb-8">
          <h2 className="text-2xl font-bold text-gray-900 mb-4 border-b pb-2">
            Education
          </h2>
          <div className="space-y-4">
            <div>
              <div className="flex justify-between items-start mb-1">
                <h3 className="text-xl font-semibold text-gray-900">
                  Degree
                </h3>
                <span className="text-gray-600">2016.09 - 2020.06</span>
              </div>
              <div className="text-gray-700">School Name</div>
            </div>
          </div>
        </section>

        <section className="mb-8">
          <h2 className="text-2xl font-bold text-gray-900 mb-4 border-b pb-2">
            Skills
          </h2>
          <div className="grid grid-cols-2 gap-4">
            <div>
              <h3 className="font-semibold text-gray-900 mb-2">Languages</h3>
              <ul className="text-gray-700 list-disc list-inside">
                <li>Go</li>
                <li>JavaScript/TypeScript</li>
                <li>Python</li>
              </ul>
            </div>
            <div>
              <h3 className="font-semibold text-gray-900 mb-2">Frameworks / Tools</h3>
              <ul className="text-gray-700 list-disc list-inside">
                <li>React</li>
                <li>Gin</li>
                <li>Docker</li>
              </ul>
            </div>
          </div>
        </section>

        <div className="print-hidden mt-8 text-center">
          <button
            onClick={() => window.print()}
            className="px-6 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 transition-colors"
          >
            Print / Export PDF
          </button>
        </div>
      </div>
    </div>
  )
}

export default Resume
