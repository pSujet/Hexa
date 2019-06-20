package Movelegtest

import (
	"math"
	"mind/core/framework/drivers/hexabody"
	"mind/core/framework/log"
	"mind/core/framework/skill"
	"mind/core/framework/drivers/distance"
	"time"
	"encoding/json"
	"strconv"
)

//Declare variable

var Control_input float64
var cpg_w [2][2]float64
var cpg_activity [2]float64
var cpg_output [2]float64
var cpg_bias float64

// CPG variable
var pcpg_step [2]float64
var set [2]float64
var setold [2]float64
var countup [2]float64
var countupold [2]float64
var countdown [2]float64
var countdownold [2]float64
var diffset [2]float64
var deltaxup [2]float64
var deltaxdown [2]float64
var xup [2]float64
var xdown [2]float64
var yup [2]float64
var ydown [2]float64
var pcpg_output0 [2]float64 //joint0
var pcpg_output1 [2]float64 //joint1

// Output to joint
var Leg_out [][]float64
var act int //0 = stop, 1 = start

// Delay
var pcpg_d_output0_0 []float64
var pcpg_d_output0_1 []float64

// Direction
var Direct float64 //True angle
var direct int //for shifting leg

// Send Data
type Message struct {
    Type string
    Data string
}
var m Message



type Movelegtest struct {
	skill.Base
}

func NewSkill() skill.Interface {
	// Use this method to create a new skill.

return &Movelegtest{}
}

// ======Control function ======

func j0out2degA(x float64) float64{
	return x*15+80 //[65-95 deg] zero at 80 deg
}

func j0out2degB(x float64) float64{
	return x*30+90 //[60-120 deg] zero at 90 deg
}

func j0out2degC(x float64) float64{
	return x*15+100 //[85-115 deg] zero at 100 deg
}

func j1out2deg(x float64) float64{
	return x*-25+60 //[50-70 deg] zero at 60 deg 
}

// ====Control function end=====

func (d *Movelegtest) OnStart() {
	hexabody.Start()
	hexabody.Stand()
	distance.Start()
	// Use this method to do something when this skill is starting.
}

func (d *Movelegtest) OnClose() {
	hexabody.StopWalkingContinuously()
	hexabody.Relax()
	distance.Close()
	hexabody.Close()

	// Use this method to do something when this skill is closing.
}

func (d *Movelegtest) OnConnect() {
	// Use this method to do something when the remote connected.
	hexabody.MoveHead(0,0)
	time.Sleep(2*time.Second)
	log.Info.Println("Start")

		//==============Control start=============
		
		//***********CPG*****************	
		Control_input = 0.169
		// Control_input = 0.02 Chaos wave
		// Control_input = 0.03 Chaos wave
		// Control_input = 0.028 slow wave gait*************
		// Control_input = 0.035 better Chaos wave can walk
		// Control_input = 0.0358 fast wave gait***************
		// Control_input = 0.04 better Chaos wave can walk
		// Control_input = 0.045 better Chaos wave can walk
		// Control_input = 0.05 Chaos tetrapod can walk
		// Control_input = 0.054 tetrapod can walk***********
		// Control_input = 0.06 better Chaos tetrapod
		// Control_input = 0.08 right caterpillar
		// Control_input = 0.088 caterpillar**************
		// Control_input = 0.09 right caterpillar
		// Control_input = 0.1 right caterpillar
		// Control_input = 0.11 right caterpillar
		// Control_input = 0.169 tripod****************

		cpg_output[0] = 0.01
		cpg_output[1] = 0.01
		cpg_w[0][0] = 1.4
		cpg_w[1][1] = 1.4
		cpg_bias = 0

		//Delay
		pcpg_d_output0_0 := make([]float64, 300)
		pcpg_d_output0_1 := make([]float64, 300)

		//Output
		Leg_out := make([][]float64, 6)
		for j := range Leg_out {
			Leg_out[j] = make([]float64, 3)
		}
		

	  go func(){ for /*ii := 1; ii <= 3000; ii++*/ {

		if act == 0 { //stop
			hexabody.MoveHead(Direct,0)
			hexabody.Stand()
			

		}else if act == 1{ //start

		cpg_w[0][1] = 0.18 + Control_input
		cpg_w[1][0] = -0.18 - Control_input
				
		cpg_activity[0] = cpg_w[0][0]*cpg_output[0]+
						cpg_w[0][1]*cpg_output[1]+cpg_bias
		cpg_activity[1] = cpg_w[1][1]*cpg_output[1]+
						cpg_w[1][0]*cpg_output[0]+cpg_bias
		for i := 0; i < len(cpg_output); i++ {
			cpg_output[i] = math.Tanh(cpg_activity[i])
		}

		//***********CPG end*************
							
		//*******Post-Processing*********
		//From CPG
		//->sawtooth(joint 0)
		//->sawtooth with zero down(joint 1)

		pcpg_step[0] = cpg_output[0]
		pcpg_step[1] = cpg_output[1]

		setold[0] = set[0]
		setold[1] = set[1]

		countupold[0] = countup[0]
		countupold[1] = countup[1]

		countdownold[0] = countdown[0]
		countdownold[1] = countdown[1]		

		// 1) Linear threshold tf step function

		if pcpg_step[0] >= 0.85{
			set[0] = 1.0
		}
		if pcpg_step[0] < 0.85{
			set[0] = -1.0
		}
		if pcpg_step[1] >= 0.85{
			set[1] = 1.0
		}
		if pcpg_step[1] < 0.85{
			set[1] = -1.0
		}
		diffset[0] = set[0] - setold[0]
		diffset[1] = set[1] - setold[1]

		// 2) Count how many steps of swing

		if set[0] == 1.0 {
			countup[0] = countup[0] + 1.0 
			countdown[0] = 0.0
		// Count how many steps of stance
		}else if set[0] == -1.0 {
			countdown[0] = countdown[0] + 1.0 
			countup[0] = 0.0
		}
		
		
		//Count how many steps of swing
		if set[1] == 1.0 {
			countup[1] = countup[1] + 1.0 
			countdown[1] = 0.0
		// Count how many steps of stance
		}else if set[1] == -1.0 {
			countdown[1] = countdown[1] + 1.0 
			countup[1] = 0.0
		}
		
		// 3) Memorized the total steps of swing and stance

		if countup[0] == 0.0 && diffset[0] == -2.0 && set[0] == -1.0{
			deltaxup[0] = countupold[0]
		}
		if countdown[0] == 0.0 && diffset[0] == 2.0 && set[0] == 1.0{
			deltaxdown[0] = countdownold[0]
		}
		if countup[1] == 0.0 && diffset[1] == -2.0 && set[1] == -1.0{
			deltaxup[1] = countupold[1]
		}
		if countdown[1] == 0.0 && diffset[1] == 2.0 && set[1] == 1.0{
			deltaxdown[1] = countdownold[1]
		}

		// 4) Compute y up and down

		xup[0] = countup[0]
		xdown[0] = countdown[0]

		xup[1] = countup[1]
		xdown[1] = countdown[1]

		yup[0] = ((2./deltaxup[0])*xup[0])-1
		ydown[0] = ((-2./deltaxdown[0])*xdown[0])+1

		yup[1] = ((2./deltaxup[1])*xup[1])-1
		ydown[1] = ((-2./deltaxdown[1])*xdown[1])+1

		// 5) Combine y up and down

		if set[0] >= 0.0{
			pcpg_output0[0] = yup[0]
			pcpg_output1[0] = yup[0]
		}
		if set[0] < 0.0{
			pcpg_output0[0] = ydown[0]
			pcpg_output1[0] = -1
		}
		if set[1] >= 0.0{
			pcpg_output0[1] = yup[1]
			pcpg_output1[1] = yup[1]
		}
		if set[1] < 0.0{
			pcpg_output0[1] = ydown[1]
			pcpg_output1[1] = -1
		}
		// Limit upper and lower limit
		if pcpg_output0[0] > 1.0{
			pcpg_output0[0] = 1.0 
			pcpg_output1[0] = 1.0 
		}else if pcpg_output0[0] < -1.0{
			pcpg_output0[0] = -1.0
		}

		if pcpg_output0[1] > 1.0{
			pcpg_output0[1] = 1.0 
			pcpg_output1[1] = 1.0 
		}else if pcpg_output0[1] < -1.0{
			pcpg_output0[1] = -1.0
		}
		//******Post-Processing end******

		//*********Delay***********
		
		var gap int //delay interval

		gap = 16
		
		//Que
		pcpg_d_output0_0 = append(pcpg_d_output0_0,pcpg_output0[0])
		pcpg_d_output0_0 = pcpg_d_output0_0[1:]		
		
		pcpg_d_output0_1 = append(pcpg_d_output0_0,pcpg_output1[0])
		pcpg_d_output0_1 = pcpg_d_output0_1[1:]

	//>>>>>>>>>>> Output to Hexa <<<<<<<<<<<<<

	    //     L           R
		//     0           1
		//       -       -
		//         - ^ -
		//   5 - - - + - - - 2
		//         - - -
        //       -       -
        //     4           3 

		//Tripod using neuron 0

		// L0_j0_out = j0out2degA(pcpg_output0[0]) 
		// L2_j0_out = j0out2degA(pcpg_output0[0])
		// L4_j0_out = j0out2degA(pcpg_output0[0]) 

		// L1_j0_out = 180-j0out2degA(pcpg_output0[0]) 
		// L3_j0_out = 180-j0out2degA(pcpg_output0[0]) 
		// L5_j0_out = 180-j0out2degA(pcpg_output0[0]) 

		// Moving Forward head at the middle between two leg (Joint 2 is fixed)
		
		var j2 float64

		j2 = 150 //fixed joint 2
		
		// --Right side--

		// Leg 1

		Leg_out[1][0] = j0out2degA(pcpg_d_output0_0[0]) //90
		Leg_out[1][1] = j1out2deg(pcpg_d_output0_1[0]) //80  
		Leg_out[1][2] = j2 	

		// Leg 2

		Leg_out[2][0] = j0out2degB(pcpg_d_output0_0[gap]) //j0out2degB(pcpg_output0[0]) 
		Leg_out[2][1] = j1out2deg(pcpg_d_output0_1[gap])  //j1out2deg(pcpg_output1[0]) 
		Leg_out[2][2] = j2
		 
		// Leg 3
			
		Leg_out[3][0] = j0out2degC(pcpg_d_output0_0[2*gap]) //j0out2degC(pcpg_output0[0]) 
		Leg_out[3][1] = j1out2deg(pcpg_d_output0_1[2*gap]) //j1out2deg(pcpg_output1[0])
		Leg_out[3][2] = j2

		// --Left side--

		// Leg 0

		Leg_out[0][0] = 200-j0out2degC(pcpg_d_output0_0[3*gap]) //90
		Leg_out[0][1] = j1out2deg(pcpg_d_output0_1[3*gap]) //80
		Leg_out[0][2] = j2

		// Leg 5
		
		Leg_out[5][0] = 180-j0out2degB(pcpg_d_output0_0[4*gap]) 
		Leg_out[5][1] = j1out2deg(pcpg_d_output0_1[4*gap]) 
		Leg_out[5][2] = j2
		
		// Leg 4
		 
		Leg_out[4][0] = 160-j0out2degA(pcpg_d_output0_0[5*gap]) 
		Leg_out[4][1] = j1out2deg(pcpg_d_output0_1[5*gap]) 
		Leg_out[4][2] = j2
	

	//==============Controlo end===============

	// Command to Hexa

		// Object avoidance

		// dist,_:= distance.Value()
		// if dist<500{
		// 	Direct = hexabody.Direction() 
		// 	Direct = Direct + 30
		// }else if dist>=500 && dist<1000{
		// 	Direct = Direct - 30
		// 	if Direct <=0{
		// 		Direct = 0
		// 	}
		// }
		// hexabody.MoveHead(Direct,0)
		
		hexabody.MoveHead(Direct,0)
		direct := int((int(Direct+30)%360)/60) //for shifting leg
				
		

		// --Right side--

		// Leg 1

		go hexabody.MoveJoint(1,0,Leg_out[(1+direct)%6][0],0)
		go hexabody.MoveJoint(1,1,Leg_out[(1+direct)%6][1],0)
		go hexabody.MoveJoint(1,2,Leg_out[(1+direct)%6][2],0)

		// // Leg 2

		go hexabody.MoveJoint(2,0,Leg_out[(2+direct)%6][0],0)
		go hexabody.MoveJoint(2,1,Leg_out[(2+direct)%6][1],0)
		go hexabody.MoveJoint(2,2,Leg_out[(2+direct)%6][2],0)

		// // Leg 3

		go hexabody.MoveJoint(3,0,Leg_out[(3+direct)%6][0],0)
		go hexabody.MoveJoint(3,1,Leg_out[(3+direct)%6][1],0)
		go hexabody.MoveJoint(3,2,Leg_out[(3+direct)%6][2],0)

		// // --Left side--

		// // Leg 0

		go hexabody.MoveJoint(0,0,Leg_out[(0+direct)%6][0],0)
		go hexabody.MoveJoint(0,1,Leg_out[(0+direct)%6][1],0)
		go hexabody.MoveJoint(0,2,Leg_out[(0+direct)%6][2],0)

		// // Leg 5

		go hexabody.MoveJoint(5,0,Leg_out[(5+direct)%6][0],0)
		go hexabody.MoveJoint(5,1,Leg_out[(5+direct)%6][1],0)
		go hexabody.MoveJoint(5,2,Leg_out[(5+direct)%6][2],0)

		// // Leg 4

		go hexabody.MoveJoint(4,0,Leg_out[(4+direct)%6][0],0)
		go hexabody.MoveJoint(4,1,Leg_out[(4+direct)%6][1],0)
		go hexabody.MoveJoint(4,2,Leg_out[(4+direct)%6][2],0)


	
	// time.Sleep(time.Second)
	// log.Info.Println("Sleep")
		// log.Info.Println(ii)
		//log.Info.Println(dist)
		//log.Info.Println(Direct)
		// log.Info.Print("Position: ",L2_j1_out)
		// log.Info.Print("Position: ",L2_j1_out)
		// log.Info.Println("Position: ",L2_j2_out)
	
	}}
}()
}

func (d *Movelegtest) OnDisconnect() {
	hexabody.Relax()
	// Use this method to do something when the remote disconnected.
}

func (d *Movelegtest) OnRecvJSON(data []byte) {
	// Use this method to do something when skill receive json data from remote client.
	json.Unmarshal(data, &m)
	Control_input,_ = strconv.ParseFloat(m.Data, 64)
	log.Info.Println(Control_input)
}

func (d *Movelegtest) OnRecvString(data string) {
	// Use this method to do something when skill receive string from remote client.
	switch data {
	case "slow_wave":
		Control_input = 0.028
	case "fast_wave":
		Control_input = 0.0358
	case "tetrapod":
		Control_input = 0.054
	case "caterpillar":
		Control_input = 0.088
	case "tripod":
		Control_input = 0.169
	case "CCW":
		Direct = Direct+30
		if Direct >= 360{
			Direct = 0
		}
	case "CW":
		Direct = Direct-30
		if Direct <= 0{
			Direct = Direct+360
		}
	case "start":
		act = 1
	case "stop":
		act = 0
	}
	log.Info.Println("act",act)
	log.Info.Println("Control Input",Control_input)
	log.Info.Println("Angle",Direct)
}
