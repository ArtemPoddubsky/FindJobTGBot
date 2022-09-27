# FindJob Telegram Bot
A simple telegram bot written in Golang which can send you vacancies using HH API and save those in history. <a href=https://t.me/InternshipGolangBot>@BotLink</a>

<h2>Overview</h2>

Бот принимает название вакансии и уточняющие параметры поиска, обращается к HeadHunter API с сформированным запросом, предоставляет ответ пользователю и сохраняет ссылки на вакансии в БД. При повторном запросе пользователю отправляются только те вакансии, которых ещё нет в БД.

Данные из БД можно очистить командой /clear.

Команда /repeat вытаскивает из БД последний отправленный запрос и отправляет его автоматически.

Список команд также доступен в меню бота, см. <a href=https://t.me/InternshipGolangBot>@BotLink</a>.

На данный момент бот будет искать только вакансии, доступные в Москве.

<h2>How to use:</h2>

      # Put your Telegram Token (given by BotFather) into configuration file /Bot/config/config.toml
  
      make all
      
      make lint - Запуск линтера
