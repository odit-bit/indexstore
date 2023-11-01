package postgreindex

// like matchDocQuery but whe $1 is blank or space it will return all
const searchDocCountQuery = `
SELECT COUNT(*) FROM documents
WHERE
	CASE
		WHEN length(trim($1)) = 0 THEN true
		ELSE ts @@ websearch_to_tsquery('english', $1)
	END
`

// plainto_tsquery vs to_tsquery
var searchDocQuery = `
SELECT linkID, url, title, content, indexed_at, pagerank
FROM documents
WHERE
	CASE
		WHEN length(trim($1)) = 0 THEN true
		ELSE ts @@ websearch_to_tsquery('english', $1)
	END
ORDER BY
	pagerank DESC,

	CASE
		WHEN length(trim($1)) = 0 THEN NULL
	 	ELSE ts_rank(ts, websearch_to_tsquery('english', $1), 32)
	END DESC	

OFFSET ($2) ROWS
FETCH FIRST ($3) ROWS ONLY;
`

const updateScoreQuery = `
	UPDATE documents
	SET pagerank = $1 -- Replace with the pagerank value
	WHERE linkID = $2; -- Replace with the specific linkID 

`

const findDocumentQuery = `
	SELECT linkID, url, title, content, indexed_at, pagerank FROM documents
	WHERE linkID = $1
`

const insertDocumentQuery = `
	INSERT INTO documents (linkID, url, title, content, indexed_at, pagerank)
	VALUES($1,$2,$3,$4, $5, $6)
	ON CONFLICT (linkID) DO 
	UPDATE
		SET url = EXCLUDED.url,
			title = EXCLUDED.title,
			content = EXCLUDED.content,
			indexed_at = NOW();
`
