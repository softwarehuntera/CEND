# CEND
Central Encyclopedia for Name Disambiguation (CEND)


# Front End
Uses tools like

https://www.npmjs.com/
https://nextjs.org/
https://ui.shadcn.com/

### API Interaction
```
// A 'category' must be specified for each operation
// E.g. we are searching over the 'Company' category, or maybe the 'Fruit' category.
// Operations provided by the backend are as follows:
view terms
delete term
add term
search term
```

# Back End
### Design decisions
- Databases will each reflect a different category.
- The concept of 'disambiguation' will be totally defined by the user. We just provide a tool that automates disambiguation recommendations.
- Functions to 'merge' and 'split' entries will be exposed, for entity disambiguation.


### The hard part
- Providing recommendations for searches.
Identify the use-case where there are two insances of the name 'Henry Matthews' inside a 'People' taxonomy; There's no way of distinguishing the two. We can search for 'Henry Matthews' and only one result will turn up. But in reality there are multiple 'Henry Matthews' floating around. How do we identify them uniquely? Based on external fields, we might need to use 'middle name' or 'email' or 'phone number'. Then with this additional information we want to disambiguate the two.

So we need to build a function that 'splits' users. What this will do is simply add information to an existing null field.
```
NAME EMAIL
Jon, None

// becomes
NAME EMAIL
Jon, jon1@gmail.com
Jon, jon2@gmail.com
```
So the contract is that we should not have duplicate rows in our database. I believe we need to be extremely careful about
defining the table schema. Tables are logically organized by category.
Each table has a unique schema, the fields are 'the minimum possible unique identifiers needed to disambiguate two entities.'

One difficult task that our application has to cover is the idea of 'names over time' - where the data that makes up a unique person 'henry matthews' changes over time. For example if he changes his email address. So then we have some metaphysical concept of what a 'henry matthews' is --- and all of his associated email addresses are linked to him.

Need to define operations like the 'merge' and 'split' of two users. Need to define the default behavior when adding duplicate entities, up to a certain field. I.e. if name and email match, then our 'address' field is added to the same user.
So maybe we have columns that are 'identifier columns' where you *must* match one of the identifiers if we are referring to the same entity. Then all extraneous information is stored in some external database. So I guess this is where an 'id' field comes in. We can use the 'id' to mark each row and then store additional metadata on the user in a table indexed by 'ids' where we know that each id refers to a completely unique entity.


### Smart Goals
- Implement basic database operations tests


### Design Decision Issues

After reading more about modern text search, we know that tools like Lucene and Elasticsearch
utilize the NoSQL database schema. The schema does not need to be specified before-hand. And an inverse
document store is used for indexing of the documents. So I need to think about my database design, with respect to search.
How do I do this? Well, we can bake in the search algorithm to the database design, and assume that everything will function on
tf-idf. This could work in a relational database, but I believe the 'finding' of documents will be much more difficult.


https://towardsdatascience.com/improving-named-entity-disambiguation-using-entity-relatedness-within-wikipedia-92f400ee5994

http://arxiv.org/pdf/1802.01074