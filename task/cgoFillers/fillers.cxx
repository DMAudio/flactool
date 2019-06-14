#include "fillers.h"
std::vector<std::vector<std::string>> GetArgs(const std::string input) {
    std::regex expr("\\{@([a-zA-Z0-9]+):(((?!(\\{@|@\\})).)*)@\\}");

    const std::sregex_token_iterator end;
    int subMatches[] = {0, 1, 2};

    std::vector<std::vector<std::string>> results;
    for (std::sregex_token_iterator itr(input.begin(), input.end(), expr, subMatches); itr != end;) {
        std::vector<std::string> result;
        result.push_back((*itr).str());
        ++itr;
        result.push_back((*itr).str());
        ++itr;
        result.push_back((*itr).str());
        ++itr;

        results.push_back(result);
    }

    return results;
}