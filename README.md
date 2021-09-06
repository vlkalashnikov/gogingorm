# gogingorm

- @Param page query string false "page"
- @Param per_page query string false "per_page"
- @Param sort query string false "Сортировка asc/desc" Enums(asc, desc)
- @Param sort_field query string false "Поля для сортировки через запятую"
- @Param filter_values query string false "Значение поля для фильтрации, через запятую"
- @Param filter_fields query string false "Название полей для фильтрации, через запятую"
- @Param filter_type query string false "Тип фильтра полное совпадение (=) или совпадение по like" Enums(=, like)
- @Param filter_logic query string false "Пересечение или дополнение" Enums(and, or)
