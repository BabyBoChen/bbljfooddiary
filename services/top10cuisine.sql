SELECT A.cuisine_id, A.cuisine_name, A.unit_price,B.cuisine_type_name
,CASE WHEN A.is_one_set = true
	THEN 'YES'
	ELSE 'NO'
	END AS is_one_set
,A.review,A.last_order_date,A.restaurant,A.address,A.remark
FROM public.cuisine AS A
LEFT JOIN public.cuisine_type AS B ON A.cuisine_type=B.cuisine_type_id
WHERE A.cuisine_type=1
ORDER BY review DESC, cuisine_name ASC
LIMIT 10

SELECT A.cuisine_id, A.cuisine_name, A.unit_price,B.cuisine_type_name
,CASE WHEN A.is_one_set = true
	THEN 'YES'
	ELSE 'NO'
	END AS is_one_set
,A.review,A.last_order_date,A.restaurant,A.address,A.remark
FROM public.cuisine AS A
LEFT JOIN public.cuisine_type AS B ON A.cuisine_type=B.cuisine_type_id
WHERE A.cuisine_type=2
ORDER BY review DESC, cuisine_name ASC
LIMIT 10

SELECT A.cuisine_id, A.cuisine_name, A.unit_price,B.cuisine_type_name
,CASE WHEN A.is_one_set = true
	THEN 'YES'
	ELSE 'NO'
	END AS is_one_set
,A.review,A.last_order_date,A.restaurant,A.address,A.remark
FROM public.cuisine AS A
LEFT JOIN public.cuisine_type AS B ON A.cuisine_type=B.cuisine_type_id
WHERE A.cuisine_type=3
ORDER BY review DESC, cuisine_name ASC
LIMIT 10