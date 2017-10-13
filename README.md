# Analytics
Analytics software was meant as a basic replacement for Google analytics
for my own website, since most of the modern Adblockers block Google analytics,
the data coming from it is probably not the actual data. All I am
interested in is the actual traffic that appears on my website and since it's hosted
on Github I can not process requests on server side, that's why this project
was made.


Do not use this software on critical high traffic site as something will probably break.

# TODO

[] Add authentication

[] Create admin dashboard for managing websites and displaying traffic data


# Database design

## Land
Land table is used for storing user traffic that appears on your website

    (id serial primary key,
    website_id integer references Website(id) on delete cascade,
    time TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    ip varchar(16));

## Click
UserClicks table is used to store user clicks on any link on your website

    (id serial primary key,
    website_id integer references Website(id) on delete cascade,
    time TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    ip varchar(16),
    url_clicked text);

## Website
Website table is used to store website urls that you want analytics for.

    (id serial primary key,
    short_url varchar(8) UNIQUE,
    owner integer references users(id) on delete cascade,
    active boolean,
    website_url text);

## Users

    (id serial primary key,
    username varchar(25) unique NOT NULL,
    email text unique NOT NULL,
    password varchar(60) NOT NULL);