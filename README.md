**Simple forum using Golang, GraphQL and PostgreSQL (optional)**

List of available queries:

getPostsPage(limit: Int = 10, offset: Int = 0): [Post]

getPostWithComments(input: ID!, limit: Int = 10, offset: Int = 0) : PostWithComments

List of available mutations:

createPost(input: NewPost!): Post!

addCommentOnPost(input: NewComment!): Comment!

addCommentOnComment(input: NewComment!): Comment!

**To build and launch docker:**

docker-compose up --build
