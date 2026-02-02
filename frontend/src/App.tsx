import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { AuthProvider } from './contexts/AuthContext'
import { ProtectedRoute } from './components/ProtectedRoute'
import Login from './pages/Login'
import Dashboard from './pages/Dashboard'
import InvitationCodes from './pages/InvitationCodes'
import UserManagement from './pages/UserManagement'
import MerchantInfo from './pages/MerchantInfo'
import ServiceProviderInfo from './pages/ServiceProviderInfo'
import CreatorProfile from './pages/CreatorProfile'

function App() {
  return (
    <AuthProvider>
      <BrowserRouter>
        <Routes>
          <Route path="/login" element={<Login />} />
          <Route
            path="/dashboard"
            element={
              <ProtectedRoute>
                <Dashboard />
              </ProtectedRoute>
            }
          />
          <Route
            path="/invitations"
            element={
              <ProtectedRoute>
                <InvitationCodes />
              </ProtectedRoute>
            }
          />
          <Route
            path="/user-management"
            element={
              <ProtectedRoute>
                <UserManagement />
              </ProtectedRoute>
            }
          />
          <Route
            path="/merchant"
            element={
              <ProtectedRoute>
                <MerchantInfo />
              </ProtectedRoute>
            }
          />
          <Route
            path="/service-provider"
            element={
              <ProtectedRoute>
                <ServiceProviderInfo />
              </ProtectedRoute>
            }
          />
          <Route
            path="/creator"
            element={
              <ProtectedRoute>
                <CreatorProfile />
              </ProtectedRoute>
            }
          />
          <Route path="/" element={<Navigate to="/dashboard" replace />} />
        </Routes>
      </BrowserRouter>
    </AuthProvider>
  )
}

export default App
