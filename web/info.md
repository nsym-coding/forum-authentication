Change the organisation of folders so that all web stuff is in the web package

Used Auth middleware to check if user is logged in before they can access the create a post page

The cookie check middleware will log out user from loginauth and create a post page if you log in from a different browser and refresh

But it won't let the new log in access the create a post page it logs them out as well without deleting their cookie

Also need to check if you log in from one client it doesn't allow you to access create post on a different client without logging in on that client

Need to rethink the cookie check logic for middleware implementation as there'll different type of cookie checks depending on the request (i.e. when they log in, log out, try to post, comment etc)

POTENTIAL MIDDLEWARE

- Check if theyre logged in so we know guest or not
- When they log in check if they're logging in for the first time - if they are give them a cookie

if not first time, check other session, delete other session&cookie and then give them a new cookie

--------------------------------------------------------------
Got the comments stored on the DB now
Still not catching userID even when using link that gets userid from the URL. 


///////////////////////////////////////////////////////
-got comments showing on the store comment page for that specific comment
-it links user id to comment when going through create comment the first time but not after that?
-need to pass through the comment data to display it on a display comments page but this probably needs to be incorporated into the show post handler
