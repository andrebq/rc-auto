firmware-upload:
	scp ./firmware/pizero/motor_ctl.py $(rcCarUser)@$(rcCarIP):motor_ctl.py

firmware-download:
	scp $(rcCarUser)@$(rcCarIP):motor_ctl.py ./firmware/pizero/motor_ctl.py

firmware-run:
	ssh $(rcCarUser)@$(rcCarIP) -- 'sudo python3 /home/pi/motor_ctl.py'