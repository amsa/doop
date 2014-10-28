#Doop Tech Spec

####Author: Wenbin Xiao, Amin Saeidi

##Intro
Doop provides git-like branches for data on database level.

##Data In Doop
Like all RDBMS, data in doop is organized by tables which contain fields and rows. 
Besides that, Doop provides another layout of abstraction: branch.

A database can have different branches such that in different branch, 
data and schema of tables can be different; also, you can merge the data/scheme from different branches into one.

To address the mechanism, assuming we already have a database, we have following notations:

* `B[i]` denotes branch `i` of database.
* `B[i].T[j]` denotes table `j` in branch `i`; this is the table can be seen by users.
* `B[i].Snapshot[j]` denotes snapshot of table `j` in branch `i`.
* `B[i].HSection[j]` denotes the `horizonal section` of table `j` in branch `i`.
* `B[i].VSection[j]` denotes the `vertical section` of table `j` in branch `i`.
* `B[i].CSection[j]` denotes the `cross section` of table `j` in branch `i`.
* `B[i].RDel[j]` denotes the deleted rows' primary keys in table `j` in branch `i`.
* `B[i].CDel[j]` denotes the dropped columns' names in table `j` in branch `i`.

`B[i].T[j]` is the table seen by user. It's a logical table which means that data is 
not physically arranged as how `B[i].T[j]` looks like; instead, data is 
stored in `B[i].Snapshot[j]`, `B[i].HSection[j]`, `B[i].VSection[j]`, `B[i],CSection[j]`, `B[i].RDel[j]`, `B[i].CDel[j]`, 
we call them `companion sections` of `B[i].T[j]`. 

Write(insert, delete, update rows, add and drop columns) operations on `B[i].T[j]` go 
to `B[i].HSection[j]` as well as other companion sections; and when there is queries on `B[i].T[j]`, 
we reassemble `B[i].T[j]` from its companion sections. 

Representing `B[i].T[j]` by graph, it looks like:

    ---------------------
    |               | V |
    |  Snapshot     | S |
    |               |   |
    |               |   |
    |               |   |
    |               |   |
    |               |   |
    |               |   |
    |               |   |
    |---------------|---|
    |               | C |
    |  HS           | S |
    |---------------|---|

    RDel: rows deleted in snapshot
    CDel: columns deleted in snapshot

###Schema and Roles of Companion Sections

* Snapshot is original table when branch is created;
* HSection has same schema as Snapshot
* VSection's schema contains newly-added columns in this branch as well as a foreign key to Snapshot; Snapshot and VSection have 1-to-1 relation.
* CSection has same schema as VSection; CSection has 1-to-1 relation to HSection and references HSection's primary key and set to `ON DELETE CASCADE`.
* RDel: the only column in RDel in the foreign key referencing primary key of Snapshot.
* CDel: only column in CDel, it's string type, contains the name of deleted columns in Snapshot. 

###Terms We Use

* Columns in Snapshot is `snapshot_columns`.
* Columns in VSection, excluding the foreign key pointing to the row in snapshot, are `new_columns`.
* Columns in `CDel`'s rows, we call them `phantom_columns`.
* Columns in `snapshot_columns` excluding `phantom_columns`, we call them `concrete_snapshot_columns`.

##Write on `B[i].T[j]`

###Patterns We can reuse

    ;this gives us all valid rows in snapshot with new columns
    ;we call it valid_upper_section
    SELECT concrete_snapshot_columns, new_columns FROM 
    snapshot 
    NATURAL INNER JOIN 
    vsection
    WHERE 
    concrete_snapshot_columns.key NOT IN
    (SELECT * FROM rdel);

    ;this gives us all valid rows in hsection with csection
    ;we call it valid_lower_section
    SELECT concrete_snapshot_columns, new_columns FROM
    hsection
    NATURAL INNER JOIN
    csection;

###Insert new records
In `B[i]`, we issue the statement: 

    INSERT INTO t VALUES (value1,value2,value3,...,value_n);

Basically, we have:

* `hsection_values`: values for `snapshot_columns`; added NULL for values for `phantom_columns`; 
* `csection_values`: values for `new_columns` and the primary key of this record.

Then we rewrite the statement into:
        
    INSERT INTO hsection VALUES (hsection_values);
    INSERT INTO csection VALUES (csection_values);

###Add new columns
In `B[i]`, we issue the statement:
    
    ALTER TABLE t ADD column_name column_type;

The statement above becomes

    ALTER TABLE vsection ADD column_name column_type;
    ALTER TABLE csection ADD column_name column_type;

###Delete rows
In `B[i]`, we issue statement:

    DELETE FROM t WHERE iampredicate;

It becomes

    ;mark rows in snapshot
    INSERT INTO rdel SELECT key FROM
    valid_upper_section                     ;see the reusable patterns above
    WHERE iampredicate;
    
    ;delete rows in hsection
    ;note that records in csection will be deleted automatically because of cascade deletion.
    DELETE FROM hsection WHERE hsection.key IN
    (
        SELECT key FROM 
        valid_lower_section
        WHERE iampredicate;
    )
    

###Update rows
In `B[i]`, we issue the statement:

    UPDATE t SET column1=value1,column2=value2,... WHERE iampredicate;
    

We rewrite the statement into:
    
    ;======= update rows in hsection and csection =======
    SELECT key INTO affected_lower_section FROM 
        valid_lower_section
    WHERE iampredicate;

    ;update hsection
    UPDATE hsection SET column1=value1,...,column_i=value_i WHERE key IN 
        (SELECT * FROM affected_lower_section);

    ;update csection
    UPDATE csection SET columni_(i+1)=value_(i+1),...,column_n=value_n WHERE key IN 
        (SELECT * FROM affected_lower_section);

    ;======= done =======

    ;======= update rows in snapshot ========
    ;find out rows affected in snapshot
    SELECT * INTO affected_upper_section
        valid_upper_section
    WHERE iampredicate;

    ;mark rows affected in snapshot
    INSERT INTO rdel
    SELECT key FROM affected_upper_section;

    ;apply the update for rows in snapshot
    UPDATE affected_upper_section SET column1=value1,column2=value2,...; 

    ;insert hsection 
    INSERT INTO hsection
    (concrete_snapshot_columns) 
    SELECT concrete_snapshot_columns FROM affected_upper_section;

    ;insert the vsection 
    INSERT INTO vsection 
    SELECT primary_key, new_columns FROM affected_upper_section;  
    ;====== done ======
    

###Drop columns
In `B[i]`, we issue statement:

    ALTER TABLE t_j DROP COLUMN column_name;

We rewrite it to:

    if column_name in vsection
        ALTER TABLE vsection DROP COLUMN column_name; 
        ALTER TABLE csection DROP COLUMN column_name; 
    else
        INSERT INTO cdel_j VALUES (column_name); 
    endif
    
##Query on `B[i].T[j]` 

In `B[i]`, given query like:

    SELECT * FROM t_j WHERE iampredicate;

###Reconstruct `B[i].T[j]`
We can reassemble `t` by:

        (valid_upper_section)
            UNION
        (valid_lower_section);

Then we can evaluate the query on `t`. 
        

##Branch Creation
When we create branch `B[i]` forking from `B[i-1]`, we turn `B[i-1]` to `snapshot`, 
let's give it a name `s1`, then we create a branch `B[i]` and a new `B[i-1]` such that:

* For all  `j >= 0 && j < B[i].T.size` , we have `B[i].Snapshot[j] = B[i-1].Snapshot[j] = s1.T[j]`
* Create companion sections for each `B[i].T[j]`  and `B[i-1].T[j]`as empty tables.

##Branch Merging

When me merge `B[i-1]` to `B[i]`: we union their corresponding companion sections. 
If duplicate keys or column names occur during union, that means conflict which requires resolution.
