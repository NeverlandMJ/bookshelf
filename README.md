# bookshelf
Bookshelf task for Lavina Tech. 
Explaination:

I wanted to use Redis to save user secret keys like a session. 
But my free Heroku account doesn't support Redis, so I had to use cache to store secret keys.
But as the cache is a tempraray momery if the server crashes all cached values will be lost. 
Hence, if the cached value is not found we will retrive it directly from database.
