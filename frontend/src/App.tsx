import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { AuthProvider } from './contexts/AuthContext'
import { ToastProvider } from './contexts/ToastContext'
import { ProtectedRoute } from './components/ProtectedRoute'
import ErrorBoundary from './components/ErrorBoundary'
import Login from './pages/Login'
import DebugAuth from './pages/DebugAuth'
import Dashboard from './pages/Dashboard'
import InvitationCodes from './pages/InvitationCodes'
import UserManagement from './pages/UserManagement'
import MerchantInfo from './pages/MerchantInfo'
import Merchants from './pages/Merchants'
import ServiceProviderInfo from './pages/ServiceProviderInfo'
import ServiceProviders from './pages/ServiceProviders'
import CreatorProfile from './pages/CreatorProfile'
import TaskHall from './pages/TaskHall'
import MyTasks from './pages/MyTasks'
import CreateCampaign from './pages/CreateCampaign'
import Recharge from './pages/Recharge'
import CreditTransactions from './pages/CreditTransactions'
import Withdrawal from './pages/Withdrawal'
import Withdrawals from './pages/Withdrawals'
import WithdrawalReview from './pages/WithdrawalReview'

function App() {
  return (
    <ErrorBoundary>
      <ToastProvider>
        <AuthProvider>
          <BrowserRouter>
            <Routes>
              <Route path="/login" element={<Login />} />
              <Route path="/debug-auth" element={<DebugAuth />} />
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
            path="/merchants"
            element={
              <ProtectedRoute>
                <Merchants />
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
            path="/service-providers"
            element={
              <ProtectedRoute>
                <ServiceProviders />
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
          <Route
            path="/task-hall"
            element={
              <ProtectedRoute>
                <TaskHall />
              </ProtectedRoute>
            }
          />
          <Route
            path="/my-tasks"
            element={
              <ProtectedRoute>
                <MyTasks />
              </ProtectedRoute>
            }
          />
          <Route
            path="/create-campaign"
            element={
              <ProtectedRoute>
                <CreateCampaign />
              </ProtectedRoute>
            }
          />
          <Route
            path="/recharge"
            element={
              <ProtectedRoute>
                <Recharge />
              </ProtectedRoute>
            }
          />
          <Route
            path="/credit-transactions"
            element={
              <ProtectedRoute>
                <CreditTransactions />
              </ProtectedRoute>
            }
          />
          <Route
            path="/withdrawal"
            element={
              <ProtectedRoute>
                <Withdrawal />
              </ProtectedRoute>
            }
          />
          <Route
            path="/withdrawals"
            element={
              <ProtectedRoute>
                <Withdrawals />
              </ProtectedRoute>
            }
          />
          <Route
            path="/withdrawal-review"
            element={
              <ProtectedRoute>
                <WithdrawalReview />
              </ProtectedRoute>
            }
          />
          <Route path="/" element={<Navigate to="/dashboard" replace />} />
            </Routes>
          </BrowserRouter>
        </AuthProvider>
      </ToastProvider>
    </ErrorBoundary>
  )
}

export default App
