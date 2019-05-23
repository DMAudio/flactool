%module cgoFilters

%{
#include "filters.h"
%}

%insert(cgo_comment_typedefs) %{
#cgo LDFLAGS: -static-libstdc++ -static-libgcc -static
%}

%include "std_vector.i"
%include "std_string.i"

%template(VectorString) std::vector<std::string>;
%template(VectorVectorString) std::vector<std::vector<std::string>>;

%include "filters.h"