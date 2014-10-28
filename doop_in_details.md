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
we call them `companison sections` of `B[i].T[j]`. 

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

##Write on `B[i].T[j]`
###Insert new records
In `B[i]`, we issue the statement: 

    INSERT INTO t_j VALUES (value1,value2,value3,...,value_n);

First of all, `value1,value2,...,value_n` can be divided into two subsets:

* `value1,value2,...,value_i` are the values for columns already in `B[i].Snapshot[j]`, we call them `snapshot values`.
* `value_(i+1),...,value_n` are values for columns not in `Snapshot[j]`, then these values must go to `VSection[j]`. We call them `vsection values`.

Assume `key` is the primary key of `T[j]` and `value1` is the value for it, the insertion becomes:
        
    INSERT INTO hsection_j VALUES (snapshot_values);
    if vsection_values is not empty
        INSERT INTO csection_j VALUES (value1, vsection_values);
    endif

###Add new columns
In `B[i]`, we issue the statement:
    
    ALTER TABLE t_j ADD column_name column_type;

Assume `key` is the primary key of `t_j` and its type is `key_type`, the statement above becomes

    ALTER TABLE vsection_j ADD column_name column_type;
    ALTER TABLE csection_j ADD column_name column_type;

###Delete rows
In `B[i]`, we issue statement:

    DELETE FROM t_j WHERE iampredicate;

it becomes

    SELECT key INTO all FROM t_j WHERE iampredicate;    
    SELECT key INTO in_h FROM all WHERE key IN
        (SELECT key FROM hsection_j);
    if in_h is not empty
        INSERT INTO rdel_j 
            SELECT key FROM all WHERE key NOT IN in_h;
        DELETE FROM vsection_j WHERE key IN in_h;       
        DELETE FROM hsection_j WHERE key IN in_h;
    else
        INSERT INTO rdel_j SELECT key FROM all;
    endif
    
Then we remove the intermediate tables.

###Update rows
In `B[i]`, we issue the statement:

    UPDATE t_j SET column1=value1,column2=value2,... WHERE iampredicate;
    

Columns in the statement can be categorized into:
    
* `column1,column2,...,column_i` which are in snapshot_j, we call them snapshot_columns.
* `column_(i+1),...,column_n` which appear in vsection_j, we call them vsection_columns.

    SELECT key INTO all FROM t_j WHERE iampredicate;
    SELECT key INTO in_h FROM all WHERE key IN
        (SELECT key FROM hsection_j);
    if in_h is not empty
        UPDATE hsection_j SET column1=value1,column2=value2,...,column_i=value_i WHERE key IN
            (SELECT key FROM in_h);
        


    

###Drop columns
In `B[i]`, we issue statement:

    ALTER TABLE t_j DROP COLUMN column_name;
We change it to:

    if column_name in vsection_j
        ALTER TABLE vsection_j DROP COLUMN column_name; 
    else
        INSERT INTO cdel_j VALUES (column_name); 
    endif
    
##Query on `B[i].T[j]` 

In `B[i]`, given query like:

    SELECT * FROM t_j WHERE iampredicate;

Columns in `t_j` can be categorized as:

* Columns in `snapshot_j` or `vsection_j`, we call them `complete_columns`.
* Columns appear in `delc_j`, we call them `phantom_columns`.
* Columns in `complete_columns` excludes `phantom_columns`, we call them `concrete_columns`.

###Reconstruct `B[i].T[j]`
We can reassemble `t_j` by:

    SELECT * INTO all_rows FROM (
        (SELECT * FROM snapshot_j WHERE snapshot_j.key NOT IN
            (SELECT key FROM hsection_j))
            UNION
        (SELECT * FROM hsection_j)
    );
        
    SELECT * INTO t_j FROM all_rows NATURAL INNER JOIN vsection_j 
    WHERE key NOT IN 
        (SELECT key FROM rdel_j);

Then we can evaluate the query on `t_j`. 

        

##Branch Creation
When we create branch `B[i]` forking from `B[i-1]`, we turn `B[i-1]` to `snapshot`, 
let's give it a name `s1`, then we create a branch `B[i]` and a new `B[i-1]` such that:

* For all  `j >= 0 && j < B[i].T.size` , we have `B[i].Snapshot[j] = B[i-1].Snapshot[j] = s1.T[j]`
* Create companion sections for each `B[i].T[j]`  and `B[i-1].T[j]`as empty tables.

##Branch Merging

When me merge `B[i-1]` to `B[i]`: we union their corresponding companion sections. 
If duplicate keys or column names occur during union, that means conflict which requires resolution.
