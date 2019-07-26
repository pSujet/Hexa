%% Neural Control Oscillator
% VISTEC Internship 2019

%% Clear
clc;
clear all;
close all;
%% Define value
Control_input = 0.169; %Modulatory Input
B1 = 0;
B2 = 0;
a1 = 0;
a2 = 0;
H1 = 0.01;
H2 = 0.01;
w11 = 1.4; %cpg_w.at(0).at(0)
w22 = 1.4; %cpg_w.at(1).at(1)
w12 = 0.18+Control_input; %cpg_w.at(0).at(1)
w21 = -0.18-Control_input; %cpg_w.at(1).at(0)

%% ==== MODULE 1 ====
% CPG
t = 800;
time = 1:t;
count = 0;
for i=1:length(time)-1
% if i >= 300 && i<=500
%     Control_input = 0.04;
% end
% if i >= 500
%     Control_input = 0.12;
% end
w12 = 0.18+Control_input; 
w21 = -0.18-Control_input;
a1(i+1) = w11*H1(i)+w12*H2(i)+B1;
a2(i+1) = w22*H2(i)+w21*H1(i)+B2;
H1(i+1) = tanh(a1(i+1));
H2(i+1) = tanh(a2(i+1));
end

%CPG Plot
figure
plot(time,H1);
hold on
plot(time,H2);
xlim([300 600]);
grid on;
xlabel("Time[steps]")
ylabel("CPG")
title("MI = 0.16")

%% ==== MODULE 2 ====
% CPG post processing
o1_step = [];
o2_step = [];
for i=1:t
    if H1(i)>=0.85
        o1_step(i)=1;
    else
        o1_step(i)=-1;
    end
    if H2(i)>=0.85
        o2_step(i)=1;
    else
        o2_step(i)=-1;
    end  
end
figure
plot(time,o1_step);
hold on
plot(time,o2_step);
xlim([300 600]);
grid on;
xlabel("Time[steps]")
ylabel("CPG")
title("MI = 0.16")


%% Sawtooh graph
sto1 = zeros(1,t);
start1 = 1;
stop1 = 1;
sto2 = zeros(1,t);
start2 = 1;
stop2 = 1;
for i = 1:t-1
    if o1_step(i)-o1_step(i+1)==-2
        stop1 = i+1;
    end
    if o1_step(i)-o1_step(i+1)==2
        start1 = i;
    end
    if stop1>start1
        ii=start1:stop1;
        m = -2/(stop1-start1);
        c = -1-m*stop1;
        sto1(1,start1:stop1) = m*ii+c;    
    end
    if start1>stop1
        ii=stop1:start1;
        m = 2/(start1-stop1);
        c = -1-m*stop1;
        sto1(1,stop1:start1) = m*ii+c;    
    end
    
    if o2_step(i)-o2_step(i+1)==-2
        stop2 = i+1;
    end
    if o2_step(i)-o2_step(i+1)==2
        start2 = i;
    end
    if stop2>start2
        ii=start2:stop2;
        m = -2/(stop2-start2);
        c = -1-m*stop2;
        sto2(1,start2:stop2) = m*ii+c;    
    end
    if start2>stop2
        ii=stop2:start2;
        m = 2/(start2-stop2);
        c = -1-m*stop2;
        sto2(1,stop2:start2) = m*ii+c;    
    end
end
figure;
plot(time,sto1);
hold on
plot(time,sto2);
xlim([300 600]);
grid on;
xlabel("Time[steps]")
ylabel("CPG post.")
title("MI = 0.04")

%% ==== MODULE 3 ====
% VRN

H15 = tanh(H1);
H16 = tanh(1) %[ones(1,10) -ones(1,10)]);
H19 = tanh(1.7246*H15+1.7246*H16-2.48285);
H20 = tanh(-1.7246*H15-1.7246*H16-2.48285);
H21 = tanh(1.7246*H15-1.7246*H16-2.48285);
H22 = tanh(-1.7246*H15+1.7246*H16-2.48285);
H27 = tanh(0.5*H19+0.5*H20-0.5*H21-0.5*H22);
%plot(time,H1);
hold on
plot(time,H27);
grid on

%% ==== MODULE 4 ====
% PSN       
I3 = 1;
H3 = tanh(-I3+1);
H4 = tanh(I3);
H5 = tanh(-5*H3+0.5*H1);
H6 = tanh(0.5*H2-5*H4);
H7 = tanh(-5*H3+0.5*H2);
H8 = tanh(0.5*H1-5*H4);
H9 = tanh(0.5*H5+0.5);
H10 = tanh(0.5*H6+0.5);
H11 = tanh(0.5*H7+0.5);
H12 = tanh(0.5*H8+0.5);
H13 = tanh(3*H9+3*H10-1.35);
H14 = tanh(3*H11+3*H12-1.35);
figure;
plot(time,H13);
hold on
plot(time,H14);



