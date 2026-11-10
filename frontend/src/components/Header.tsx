import React from 'react'
import { ConnectionMode } from '../types'

interface HeaderProps {
  mode: ConnectionMode
  onModeChange: (mode: ConnectionMode) => void
  onSearch: (query: string) => void
  searchPlaceholder?: string
}

export const Header: React.FC<HeaderProps> = ({
  mode,
  onModeChange,
  onSearch,
  searchPlaceholder = 'Search nodes, protocols, paradigms...'
}) => {
  const [searchQuery, setSearchQuery] = React.useState('')

  const handleSearchChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value
    setSearchQuery(value)
    onSearch(value)
  }

  const handleSearchSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    onSearch(searchQuery)
  }

  return (
    <header className="header" role="banner">
      <div className="header-left">
        <h1 className="logo" aria-label="CommLayers Home">
          <svg width="28" height="28" viewBox="0 0 28 28" fill="none" aria-hidden="true" className="logo-icon">
            <circle cx="14" cy="14" r="12" stroke="currentColor" strokeWidth="2"/>
            <path d="M14 6V22M6 14H22" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" opacity="0.5"/>
            <circle cx="14" cy="14" r="5" stroke="currentColor" strokeWidth="2" fill="currentColor"/>
          </svg>
          <span className="logo-text">CommLayers</span>
        </h1>
      </div>

      <div className="header-center">
        <div className="mode-toggle" role="tablist" aria-label="Connection mode">
          <button
            role="tab"
            aria-selected={mode === 'technical'}
            aria-controls="technical-panel"
            id="technical-tab"
            className={`mode-toggle-option ${mode === 'technical' ? 'active' : ''}`}
            onClick={() => onModeChange('technical')}
          >
            Technical Connections Mode
          </button>
          <button
            role="tab"
            aria-selected={mode === 'paradigms'}
            aria-controls="paradigms-panel"
            id="paradigms-tab"
            className={`mode-toggle-option ${mode === 'paradigms' ? 'active' : ''}`}
            onClick={() => onModeChange('paradigms')}
          >
            Paradigms Mode
          </button>
        </div>
      </div>

      <div className="header-right">
        <form className="search-form" onSubmit={handleSearchSubmit} role="search">
          <label htmlFor="global-search" className="visually-hidden">Search</label>
          <div className="search-input-wrapper">
            <svg className="search-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" aria-hidden="true">
              <circle cx="11" cy="11" r="8"/>
              <path d="M21 21l-4.35-4.35"/>
            </svg>
            <input
              id="global-search"
              type="search"
              className="search-input"
              placeholder={searchPlaceholder}
              value={searchQuery}
              onChange={handleSearchChange}
              autoComplete="off"
              spellCheck={false}
            />
            {searchQuery && (
              <button
                type="button"
                className="search-clear"
                onClick={() => {
                  setSearchQuery('')
                  onSearch('')
                }}
                aria-label="Clear search"
              >
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" aria-hidden="true">
                  <line x1="18" y1="6" x2="6" y2="18"/>
                  <line x1="6" y1="6" x2="18" y2="18"/>
                </svg>
              </button>
            )}
          </div>
        </form>

        <button className="help-button" aria-label="Help and documentation" title="Help">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" aria-hidden="true">
            <circle cx="12" cy="12" r="10"/>
            <path d="M12 16v-4M12 8h.01"/>
          </svg>
        </button>

        <div className="user-avatar" aria-label="User menu" role="button" tabIndex={0}>
          <svg width="28" height="28" viewBox="0 0 28 28" fill="none" aria-hidden="true">
            <circle cx="14" cy="14" r="13" stroke="#333333" strokeWidth="1"/>
            <circle cx="14" cy="10" r="4" fill="#333333"/>
            <path d="M4 24c0-5.5 4.5-10 10-10s10 4.5 10 10" stroke="#333333" strokeWidth="1.5" fill="none"/>
          </svg>
        </div>
      </div>
    </header>
  )
}