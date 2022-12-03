import './App.css';
import CardProps from './components/Card';
import Loading from './components/Loading';

import { useState, useEffect } from 'react';

const SEARCH_ENDPOINT = process.env.REACT_APP_SEARCH_ENDPOINT
function App() {
  const [isLoading, setIsLoading] = useState(false);
  const [posts, setPosts] = useState([]);
  useEffect(() => {
      setIsLoading(true);
      fetch(SEARCH_ENDPOINT)
        .then((res) => res.json())
        .then((data) => {
            setPosts(data);
            setIsLoading(false);
        })
        .catch((err) => {
            console.log(err.message);
            setIsLoading(false);
        });
  }, []);
  return (
    <div className="App">
      <header className="App-header">
          <div className="contents">
            { isLoading ? (
              <div className="blog-container">
                <Loading />
              </div>
            ) : (
            <div className="blog-container">
              {posts.map((blog, index) => 
                <CardProps 
                  key={index}
                  Title = {blog.Title}
                  Date = {blog.Date}
                  Url = {blog.Url}
                  Company = {blog.Company}
                />
              )}
            </div>
          )}
          </div>
      </header>
    </div>
  );
}

export default App;
