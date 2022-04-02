# Enable all rules by default
all

# Extend line length, since each sentence should be on a separate line.
rule 'MD013', :line_length => 99999

# Allow in-line HTML
exclude_rule 'MD002'
exclude_rule 'MD007'
exclude_rule 'MD029'
exclude_rule 'MD031'
exclude_rule 'MD033'
exclude_rule 'MD036'
exclude_rule 'MD041'

# Nested lists should be indented with four spaces.
rule 'MD007', :indent => 4