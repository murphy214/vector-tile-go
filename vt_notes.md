# Notes on parts of v3 I've implemented


# What needs done 
    * Tag Reader / Writer 
    * Reading Tags / writing tags in tile reader / writer
    * 3d updates to geometry 
    * All the feature level stuff


# What has been done
    - 4.4.2.2 Complex Value Encoding
    - 4.4.3 Attribute Keys

# Algebra and stuff 

```
value = base + multiplier * (delta_encoded_value + offset) 
value - base = multiplier * (delta_encoded_value + offset)

delta_encoded_value  = ((value - base) / multiplier) - offset 
```

### Geometric Attributes Comments

```
4.4.2.1 Attributes and geometric attributes
The attributes of a feature are divided between the attributes and geometric_attributes messages, both of which encode a list of key-value pairs.

Feature attributes that describe the overall feature should be encoded in the attributes message.

Feature attributes that describe additional characteristics of specific locations along the feature's geometry should be encoded in the geometric_attributes message.

Each key-value pair in the geometric_attributes MUST have a value whose type is list or delta-encoded list, and whose length is the total number of moveto, lineto, and closepath commands in the geometry. Each element in the list is considered to be associated with the corresponding command in the geometry. Note that the geometric_attributes message does include data for closepath operations, unlike elevation, which does not.

```

### What the cursor currently supports 

* Hopefully geometric attributes reduction as well as slick handling of both closed polygons and unclosed polygons
* Hopefully geometric attributes reduction for elevation ignoring the close path coordinate as stated by the spec
* Hopefully correct elevation encoding but it very well could be implemented incorrectly 
* For all the MakeGeometryFloat objects, not quite sure how to handle raw coordinates yet

### What the cursor currently doesn't support

* Proper management of elevation scaling, and other potential feature level configuration variables 
* Other stuff TM 

