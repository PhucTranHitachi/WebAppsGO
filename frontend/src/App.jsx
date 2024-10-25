import {
  BrowserRouter,
  Routes,
  Route
} from 'react-router-dom'
import Home from "./pages/Home"
import AddBook from './pages/AddBook'
import UpdateBook from './pages/UpdateBook'

function App() {
  return (
    <>
      <BrowserRouter>
        <Routes>
          <Route path='/' element={<Home/>}/>
          <Route path='/add' element={<AddBook/>}/>
          <Route path='/update/:id' element={<UpdateBook/>}/>
        </Routes>
      </BrowserRouter>
    </>
  )
}

export default App
