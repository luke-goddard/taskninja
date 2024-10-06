// ========================================================
// Examples:
// ========================================================
// 1.   task add "Buy milk" due:2016-01-02 priority:high
// 2.   task 1 modify due:2016-01-02
// 3.   task list (project:home and priority:high)
// ========================================================
// Grammar:
// ========================================================
// PROGRAM -> (
// COMMAND -> (
//  COMMAND_ADD | // e.g add "Buy milk"
// )
// COMMAND_ADD -> add PARAM EXPRESSION_STATEMENTS | add EXPRESSION_STATEMENTS PARAM | add PARAM | add EXPRESSION_STATEMENTS PARAM EXPRESSION_STATEMENTS
// EXPRESSION_STATEMENTS -> EXPRESSION_STATEMENT | EXPRESSION_STATEMENT EXPRESSION_STATEMENTS
// EXPRESSION_STATEMENT -> EXPRESSION | EXPRESSION EXPRESSION_STATEMENT
// EXPRESSION -> BINARY_EXPRESSION | LOGICAL_EXPRESSION | TAG | PAIR
// BINARY_EXPRESSION -> (EXPRESSION) BINARY_OPERATOR (EXPRESSION) | (EXPRESSION) BINARY_OPERATOR (EXPRESSION) BINARY_EXPRESSION
// LOGICAL_EXPRESSION -> (EXPRESSION) LOGICAL_OPERATOR (EXPRESSION) | (EXPRESSION) LOGICAL_OPERATOR (EXPRESSION) LOGICAL_EXPRESSION
// LOGICAL_OPERATOR -> and | or
// BINARY_OPERATOR -> + | - | * | / | %
// PARAM -> TASKID | STRING
// TAG -> +TAG | -TAG
// PAIR -> key:EXPRESSION | key:EXPRESSION_STATEMENTS
// TASKID -> number

package parser
