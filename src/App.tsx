import './styles/App.css'
import Post from './Post'

function App() {

  return (
    <div className='app'>
      <div className='app__header'>
        <img
          className='app__headerImage'
          src='https://www.instagram.com/static/images/web/mobile_nav_type_logo.png/735145cfe0a4.png'
          alt=''
        />
        <Post />
      </div>
    </div>
  )
}

export default App
