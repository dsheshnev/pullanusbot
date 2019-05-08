package main

import (
	"fmt"
	"runtime"
)

var ru = map[string]string{
	"faggot_rules": `Правила игры *Пидор Дня* (только для групповых чатов):
	*1*. Зарегистрируйтесь в игру по команде /pidoreg
	*2*. Подождите пока зарегиструются все (или большинство :)
	*3*. Запустите розыгрыш по команде /pidor
	*4*. Просмотр статистики канала по команде /pidorstats, /pidorall
	*5*. Личная статистика по команде /pidorme
	*6*. Статистика за 2018 год по комнаде /pidor2018 (так же есть за 2016-2017)
	
	*Важно*, розыгрыш проходит только *раз в день*, повторная команда выведет *результат* игры.
	
	Сброс розыгрыша происходит каждый день в 12 часов ночи по UTC+2 (или два часа ночи по Москве).`,

	"faggot_not_available_for_private": "Извините, данная команда недоступна в личных чатах.",
	"faggot_added_to_game":             "Ты в игре!",
	"faggot_already_in_game":           "Эй! Ты уже в игре!",
	"faggot_no_players":                "Зарегистрированных в игру еще нет, а значит *пидор* ты - %s",
	"faggot_not_enough_players":        "Нужно как минимум два игрока, чтобы начать игру! Зарегистрируйся используя /pidoreg",
	"faggot_winner_known":              "Согласно моей информации, по результатам сегодняшнего розыгрыша *пидор дня* - %s!",
	// 0
	"faggot_game_0_0": "Осторожно! *Пидор дня* активирован!",
	"faggot_game_0_1": "Система взломана. Нанесён урон. Запущено планирование контрмер.",
	"faggot_game_0_2": "Сейчас поколдуем...",
	"faggot_game_0_3": "Инициирую поиск *пидора дня*...",
	"faggot_game_0_4": "Итак... кто же сегодня *пидор дня*?",
	"faggot_game_0_5": "Кто сегодня счастливчик?",
	"faggot_game_0_6": "Зачем вы меня разбудили...",
	"faggot_game_0_7": "### RUNNING 'TYPIDOR.SH'...",
	"faggot_game_0_8": "Woop-woop! That's the sound of da pidor-police!",
	"faggot_game_0_9": "Опять в эти ваши игрульки играете? Ну ладно...",
	// 1
	"faggot_game_1_0": "_Шаманим-шаманим_...",
	"faggot_game_1_1": "_Где-же он_...",
	"faggot_game_1_2": "_Сканирую_...",
	"faggot_game_1_3": "_Военный спутник запущен, коды доступа внутри_...",
	"faggot_game_1_4": "_Хм_...",
	"faggot_game_1_5": "_Интересно_...",
	"faggot_game_1_6": "_Ведётся поиск в базе данных_...",
	"faggot_game_1_7": "_Машины выехали_",
	"faggot_game_1_8": "_(Ворчит) А могли бы на работе делом заниматься_",
	"faggot_game_1_9": "_Выезжаю на место..._",
	// 2
	"faggot_game_2_0": "Так-так, что же тут у нас...",
	"faggot_game_2_1": "КЕК!",
	"faggot_game_2_2": "Доступ получен. Аннулирование протокола.",
	"faggot_game_2_3": "Проверяю данные...",
	"faggot_game_2_4": "Ох...",
	"faggot_game_2_5": "Высокий приоритет мобильному юниту.",
	"faggot_game_2_6": "Ведётся захват подозреваемого...",
	"faggot_game_2_7": "Что с нами стало...",
	"faggot_game_2_8": "Сонно смотрит на бумаги",
	"faggot_game_2_9": "В этом совершенно нет смысла...",
	// 3
	"faggot_game_3_0": "Ого, вы посмотрите только! А *пидор дня* то - %s",
	"faggot_game_3_1": "Кажется, *пидор дня* - %s",
	"faggot_game_3_2": ` ​ .∧＿∧
	( ･ω･｡)つ━☆・\*。
	⊂  ノ    ・゜+.
	しーＪ   °。+ \*´¨)
			 .· ´¸.·\*´¨)
			  (¸.·´ (¸.·"\* ☆ ВЖУХ И ТЫ ПИДОР, %s`,
	"faggot_game_3_3": "И прекрасный человек дня сегодня... а нет, ошибка, всего-лишь *пидор* - %s",
	"faggot_game_3_4": "Анализ завершен. Ты *пидор*, %s",
	"faggot_game_3_5": "Ага! Поздравляю! Сегодня ты *пидор*, %s",
	"faggot_game_3_6": "Что? Где? Когда? А ты *пидор дня* - %s",
	"faggot_game_3_7": "Ну ты и *пидор*, %s",
	"faggot_game_3_8": "Кто бы мог подумать, но *пидор дня* - %s",
	"faggot_game_3_9": "Стоять! Не двигаться! Вы объявлены *пидором дня*, %s",

	"faggot_stats_top":    "Топ-10 *пидоров* за текущий год:",
	"faggot_stats_entry":  "*%d*. %s — _%d раз(а)_",
	"faggot_stats_bottom": "Всего участников — _%d_",

	"faggot_all_top":    "Топ-10 *пидоров* за всё время:",
	"faggot_all_entry":  "*%d*. %s — _%d раз(а)_",
	"faggot_all_bottom": "Всего участников — _%d_",

	"faggot_me": "%s, ты был(а) *пидором дня* — %d раз!",
}

func i18n(key string) string {

	if val, ok := ru[key]; ok {
		return val
	}

	// log.Fatal(fmt.Sprintf("KEY_MISSED:\"%s\"", key))

	_, file, line, _ := runtime.Caller(0)
	return fmt.Sprintf("%s:%d KEY_MISSED:\"%s\"", file, line, key)
}

// func pluralize(key string, n int, forms ...string) string {
// 	_, file, line, _ := runtime.Caller(0)
// 	return fmt.Sprintf("%s:%d IMPLEMENT_ME:\"%s\"", file, line, key)
// }
