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
//  COMMAND_ADD |                              // e.g add "Buy milk"
//  COMMAND_LIST |                             // e.g list
// )
// COMMAND_ADD -> add PARAM OPTIONS | add OPTIONS PARAM | add PARAM | add OPTIONS PARAM OPTIONS
// COMMAND_LIST -> list OPTIONS | list
// PARAM -> TASKID | STRING
// OPTIONS -> OPTION | OPTION OPTIONS
// OPTION -> TAG | PAIR | BINARY_EXPRESSION
// BINARY_EXPRESSION -> OPTION LOGICAL_OPERATOR OPTION
// LOGICAL_OPERATOR -> and | or
// TAG -> +TAG | -TAG
// PAIR -> key:value
// TASKID -> number

package parser
