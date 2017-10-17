# Analytics
Analytics software was meant as a basic replacement for Google Analytics
for any website where you do not have the server access (ex. Github pages),
so you can not process traffic data on the server side. Most of the modern
Adblockers are blocking Google Analytics so the data coming from it is
probably not the actual data, which is why this software was made.


Note: Do not use this software on critical high traffic site as your
database will probably melt and the software will most likely change
a lot in the future.

# Screenshots

## Analytics page
![charts example](/images/charts.png)

## Dashboard
![dashboard](/images/dashboard.png)

## Edit existing website
![edit website](/images/editWebsite.png)

# TODO

- [x] Add authentication

- [x] Create admin dashboard for managing websites and displaying traffic data

- [ ] Browser version logging

- [ ] Most clicked content



# Database design

## Browser
Storing browser versions into database, so we can reference them in other
tables by their ids

    (id serial primary key,
     browser text);

## Land
Land table is used for storing user traffic that appears on your website

    (id serial primary key,
    website_id integer references Website(id) on delete cascade,
    time TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    browser_id integer references browser(id),
    ip varchar(16));

## Click
UserClicks table is used to store user clicks on any link on your website

    (id serial primary key,
    website_id integer references Website(id) on delete cascade,
    time TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    browser_id integer references browser(id),
    ip varchar(16),
    url_clicked text);

## Website
Website table is used to store website urls that you want analytics for.

    (id serial primary key,
    short_url varchar(8) UNIQUE,
    owner integer references users(id) on delete cascade,
    name varchar(40),
    active boolean,
    website_url text);

## Users

    (id serial primary key,
    username varchar(25) unique NOT NULL,
    email text unique NOT NULL,
    password varchar(60) NOT NULL,
    active bool);