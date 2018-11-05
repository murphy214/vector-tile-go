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