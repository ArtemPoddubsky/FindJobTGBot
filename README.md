# HeadHunter Telegram Bot
A simple telegram bot written in Golang which can send you vacancies using HH API and save those in history. <a href=https://t.me/InternshipGolangBot>@BotLink</a>

<h2>Overview</h2>

Идея написать этого бота возникла, когда я начал мониторить вакансии на headhunter и обнаружил, что мне приходится вручную отслеживать список новых, не просмотренных вакансий по определенным критериям.

Я решил упростить процесс поиска работы, реализовав бота, который может искать вакансии по заданному набору параметров, вести историю просмотренных вакансий и таким образом присылать мне список только новых вакансий, которые я еще не видел.

Список команд бота доступен в меню, см. <a href=https://t.me/InternshipGolangBot>@BotLink</a>.

На данный момент бот будет искать только вакансии, доступные в Москве.

<h2>How to use:</h2>

      # Put your Telegram Token (given by BotFather) into configuration file /Bot/config/config.toml
  
      make all
