import React, { useEffect, useState } from 'react'
import { SplitView } from './components/SplitView'
import './App.css'

const App: React.FC = () => {
  const [isLoaded, setIsLoaded] = useState(false)

  useEffect(() => {
    const timer = setTimeout(() => {
      setIsLoaded(true)
    }, 300)

    return () => clearTimeout(timer)
  }, [])

  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === 'Escape') {
        const event = new CustomEvent('close-detail')
        window.dispatchEvent(event)
      }
    }

    window.addEventListener('keydown', handleKeyDown)
    return () => window.removeEventListener('keydown', handleKeyDown)
  }, [])

  if (!isLoaded) {
    return (
      <div className="loading-screen">
        <div className="loading-spinner" />
        <p>Loading CommLayers...</p>
      </div>
    )
  }

  return (
    <div className="app">
      <SplitView initialDividerPosition={0.4} />
    </div>
  )
}

export default App