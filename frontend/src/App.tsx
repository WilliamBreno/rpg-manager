import { BrowserRouter, Routes, Route } from 'react-router-dom'
import PrivateRoute from './components/PrivateRoute'
import Home from './pages/Home'
import Login from './pages/Login'
import Register from './pages/Register'
import CharacterList from './pages/CharacterList'
import CharacterCreate from './pages/CharacterCreate'
import CharacterDetail from './pages/CharacterDetail'
import CharacterEdit from './pages/CharacterEdit'

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />
        <Route path="/" element={<PrivateRoute><Home /></PrivateRoute>} />
        <Route path="/characters" element={<PrivateRoute><CharacterList /></PrivateRoute>} />
        <Route path="/characters/new" element={<PrivateRoute><CharacterCreate /></PrivateRoute>} />
        <Route path="/characters/:id" element={<PrivateRoute><CharacterDetail /></PrivateRoute>} />
        <Route path="/characters/:id/edit" element={<PrivateRoute><CharacterEdit /></PrivateRoute>} />
      </Routes>
    </BrowserRouter>
  )
}

export default App