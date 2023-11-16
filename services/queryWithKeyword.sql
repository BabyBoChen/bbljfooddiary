with VI_Cuisine as 
(
    SELECT A.cuisine_id, A.cuisine_name, A.unit_price,B.cuisine_type_name
    ,CASE WHEN A.is_one_set = true
        THEN 'YES'
        ELSE 'NO'
        END AS is_one_set
    ,A.review,A.last_order_date,A.restaurant,A.address,A.remark
    FROM public.cuisine AS A
    LEFT JOIN public.cuisine_type AS B ON A.cuisine_type=B.cuisine_type_id
)
select *
from VI_Cuisine
where cuisine_name like $1 OR cuisine_type_name like $1 OR cast(last_order_date AS VARCHAR(10)) like $1 OR restaurant like $1 OR address like $1 OR remark like $1
ORDER BY last_order_date desc