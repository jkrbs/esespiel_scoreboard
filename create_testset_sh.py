import random
commands_characters = ['#!/bin/sh']
character_template = 'curl -X POST localhost:8888/api/user/create -d "name={name}&description={description}&password={password}&eigenschaft={feature}&vorliebe={preference}"'

import pandas as pd
keys = []

task_keys = {}

with open('characters.csv') as file:
	for line in file.read().splitlines():
		name, preference, feature, _, _ = line.split(';')
		print(name, preference, feature)
		pwd = random.randint(1000, 10000)
		commands_characters.append(character_template.format(name=name, preference=preference, description='', feature=feature, password=pwd))	
		keys.append((name, pwd))

with open('commands_characters.sh','w') as out:
	out.write('\n'.join(commands_characters))
df_keys = pd.DataFrame(keys)
df_keys.to_csv('char_keys.csv',index=False, header=False)


story_template = 'curl -X POST localhost:8888/api/task/create -d "title={task}&description={description}&key={key}&points={points}&storyline={storyline}"'

questlines = '''Gremiensemester   & FSR   & ascii   & HS   & PA 
Nachtwanderung    & StuWe & ascii   & CD      & ZIH
Sportler          & USZ   & USZ     & HS   & KK
Russischer Hacker & HB    & HB      & SLUB    & HS
Regelstudienzeit  & PA    & JV      & PA      & HS
Partygänger       & SLUB  & CD      & ascii   & VL
Wirtschaftsjob    & JV  & JV      & PA      & AA
Student. Hilfskraft & PA  & JV      & ZIH     & JV
Auslandssemester   & PA   & AA     & PA     & HS
Schummler         & SLUB  & VL      & HB      & PA
Große Liebe       & CD    & ascii   & ascii   & ZIH
Wohnungswechsel       & StuWe & BA      & MA      & ZIH
Allgelehrt        & VL    & ascii   & ZIH     & VL
'''.splitlines()

story_commands = ['#!/bin/sh']
for questline in questlines:
	questline = [l.strip() for l in questline.split('&')]
	storyline = questline[0]
	task_keys[storyline] = [storyline]
	for i, task in enumerate(questline[1:]):
		pwd = random.randint(1000, 10000)
		points = random.randint(10, 100)
		story_commands.append(story_template.format(task=task, description='', key=pwd, points=points, storyline=storyline))
		task_keys[storyline].append(f'{i}:{task}:{pwd}')

df_keys_tasks = pd.DataFrame(task_keys).transpose()
df_keys_tasks.to_csv('task_keys.csv',index=False, header=False)

with open('commands_story.sh','w') as out:
	out.write('\n'.join(story_commands))