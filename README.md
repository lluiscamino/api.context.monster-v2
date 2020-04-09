# api.context.monster/v2

This is the main repository for the context.monster API. This API handles the insertion of new images and keywords by
the [Archillect Context bot](https://github.com/lluiscamino/archillect-context-bot), and also provides endpoints for
searching images and keywords stored in the context.monster database.

The [context.monster](https://context.monster/) web application uses this endpoints
to display the images sent by the bot.

## API documentation

### Authentication
To use any of these routes, you first need an API Key. This key needs to be sent
in every request as a Bearer Token.

### Entities

#### Tweet
A tweet is an image analyzed by the [bot](https://github.com/lluiscamino/archillect-context-bot) and has the following
properties:
* ID: An integer used for identifying every image
* Title: Image title
* Image: Image URL
* Archillect Tweet: Tweet ID for the [@archillect](https://twitter.com/archillect) post in which the image was published
* Date: When was the image analyzed
* Pages: An array of web pages in which the same image has been found
* Matches: An array of equal images
* Partial matches: An array of similar images
* Number of keywords: Amount of keywords for the image
* Is First: Boolean that says if the image was the first image analyzed by the bot
* Is Last: Boolean that says if the image is the last image analyzed by the bot
* Ratings: An array of Ratings (described below)

##### Endpoints
1. Get Tweet
    * Method: **GET**
    * URL: ``http://api.context.monster/v2/tweets/{id}``
    * Params:
        * {id}: Integer
 2. Get Tweets
    * Method: **GET**
    * URL: ``http://api.context.monster/v2/tweets?limit={limit}&ratings={include_ratings}&order={order}``
    * Params:
        * Limit: Integer _(Optional)_
        * Ratings: Boolean _(Optional)_
        * Order: String (ASC or DESC) _(Optional)_
3. Search Tweets
    * Method: **GET**
    * URL: ``http://api.context.monster/v2/tweets/search/{needle}?limit={limit}&order={order}&ratings={include_ratings}``
    * Params:
        * Needle: Integer
        * Limit: Integer _(Optional)_
        * Order: String (ASC or DESC) _(Optional)_
        * Ratings: Boolean
4. Create Tweet
    * Method: **POST**
    * URL: ``http://api.context.monster/v2/tweets``
    * Body: 
        ```json
        {
            "image": "https://example.com/image.png",
            "archillect_tweet": 1,
            "pages": [
                "https://example1.com",
                "https://example2.com"
            ],
            "matches": [
                "https://example.com/match1.png",
                "https://example.com/match2.png",
                "https://example.com/match3.png"
            ],
            "partial_matches": [
            	"https://example.com/partial_match1.png"
        	],
            "ratings": [
                {
                    "keyword_text": "Architecture",
                    "rate": 0.9433909058570862
                },
                {
                    "keyword_text": "Architect",
                    "rate": 0.7197999954223633
                },
                {
                    "keyword_text": "Interior Design Services",
                    "rate": 0.5656846165657043
                }
            ]
        }
       ```
      
#### Keyword
A keyword is a string used to describe one or many images (tweets). It has the following properties:
* Text: The keyword itself.
* Counter: Number of times that the keyword has appeared to describe an image.
* Searches: Number of times that the keyword has been searched.
* Ratings: An array of ratings

#####Endpoints
1. Get Keyword
    * Method: **GET**
    * URL: ``http://api.context.monster/v2/keywords/{keyword}``
    * Params
        * Keyword: String
2. Get Keywords
    * Method: **GET**
    * URL: ``http://api.context.monster/v2/keywords?limit={limit}&ratings={include_ratings}&order={order}``
    * Params
        * Limit: Integer  _(Optional)_
        * Ratings: Boolean  _(Optional)_
        * Order: String (counter or searches)
3. Search Keywords
    * Method: **GET**
    * URL: ``http://api.context.monster/v2/keywords/search/{needle}?limit={limit}&ratings={include_ratings}``
    * Params
        * Needle: String
        * Limit: Integer _(Optional)_
        * Ratings: Boolean  _(Optional)_
        
#### Rating
A rating joins a tweet with a keyword and assigns a rate to this union. It has the following properties:
* Tweet
* Keyword
* Rate

#### API Key
API Keys are needed to use the API. Each key has an access level (nothing, read-only or read-write) and every request
done with a valid API Key is logged.

##### Endpoints
1. Create API Key _(Currently disabled)_
    * Method: **GET**
        * URL: ``http://api.context.monster/v2/apikeys/new?name=name``
        * Params
            * Name: String