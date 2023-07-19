-- Create Product
/*INSERT INTO products (title, description, price)
VALUES ('product_test', 'desc_test',120);*/

-- Select Product
/*SELECT p.id, p.title
FROM products p
WHERE p.title LIKE '%Coffee%'*/

-- Left Join
/*SELECT u.id, u.email, r.title as role
FROM users u
         LEFT JOIN roles r on r.id = u.role_id;*/

-- Right Join
/*SELECT p.id, p.title, c.title as category
FROM products p
         RIGHT JOIN products_categories pc on pc.product_id = p.id
         RIGHT JOIN categories c on c.id = pc.category_id;*/

-- jsonb
/*SELECT p.id,
       p.title,
       (SELECT to_jsonb(t)
        FROM (SELECT i.id, i.filename, i.url
              FROM images i
              WHERE i.product_id = p.id) as t
        LIMIT 1) as image
FROM products p;*/

-- json array
/*SELECT p.id,
       p.title,
       (SELECT COALESCE(array_to_json(array_agg(t)),'[]'::json)
        FROM (SELECT i.id, i.filename, i.url
              FROM images i
              WHERE i.product_id = p.id) as t) as images
FROM products p;*/

-- json
SELECT json_build_object(
               'id', p.id,
               'title', p.title
           ) as json
FROM products p;