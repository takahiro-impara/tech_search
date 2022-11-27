import './App.css';
import Card from './components/Card';

import { useState, useEffect } from 'react';

const endpoint = "http://teck-search-backend.staging.udacity.impara8.com:8080/techsearch/v1/blogs"

function App() {
  const [posts, setPosts] = useState([]);
  useEffect(() => {
      fetch(endpoint)
        .then((res) => res.json())
        .then((data) => {
            setPosts(data);
        })
        .catch((err) => {
            console.log(err.message);
        });
  }, []);
  return (
    <div className="App">
      <header className="App-header">
          <div className="contents">
            <div className="blog-container">
              {posts.map((blog, index) => 
                  <Card 
                    key={index}
                    Title = {blog.Title}
                    Date = {blog.Date}
                    Url = {blog.Url}
                    Company = {blog.Company}
                  />
              )}
            </div>
          </div>
      </header>
    </div>
  );
}

export default App;
