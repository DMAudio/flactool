%module cgoFilters

%{
#include "filters.h"
%}

%include "std_vector.i"
%include "std_string.i"

%template(VectorString) std::vector< std::string>;
%template(VectorVectorString) std::vector< std::vector< std::string> >;

%include "filters.h"