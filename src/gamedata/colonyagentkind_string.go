// Code generated by "stringer -type=ColonyAgentKind -trimprefix=Agent"; DO NOT EDIT.

package gamedata

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[AgentWorker-0]
	_ = x[AgentMilitia-1]
	_ = x[AgentFreighter-2]
	_ = x[AgentRedminer-3]
	_ = x[AgentCrippler-4]
	_ = x[AgentFighter-5]
	_ = x[AgentScavenger-6]
	_ = x[AgentCourier-7]
	_ = x[AgentTrucker-8]
	_ = x[AgentPrism-9]
	_ = x[AgentServo-10]
	_ = x[AgentRepeller-11]
	_ = x[AgentDisintegrator-12]
	_ = x[AgentRepair-13]
	_ = x[AgentCloner-14]
	_ = x[AgentRecharger-15]
	_ = x[AgentGenerator-16]
	_ = x[AgentMortar-17]
	_ = x[AgentAntiAir-18]
	_ = x[AgentRefresher-19]
	_ = x[AgentStormbringer-20]
	_ = x[AgentDestroyer-21]
	_ = x[AgentMarauder-22]
	_ = x[AgentKindNum-23]
	_ = x[AgentGunpoint-24]
}

const _ColonyAgentKind_name = "WorkerMilitiaFreighterRedminerCripplerFighterScavengerCourierTruckerPrismServoRepellerDisintegratorRepairClonerRechargerGeneratorMortarAntiAirRefresherStormbringerDestroyerMarauderKindNumGunpoint"

var _ColonyAgentKind_index = [...]uint8{0, 6, 13, 22, 30, 38, 45, 54, 61, 68, 73, 78, 86, 99, 105, 111, 120, 129, 135, 142, 151, 163, 172, 180, 187, 195}

func (i ColonyAgentKind) String() string {
	if i >= ColonyAgentKind(len(_ColonyAgentKind_index)-1) {
		return "ColonyAgentKind(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ColonyAgentKind_name[_ColonyAgentKind_index[i]:_ColonyAgentKind_index[i+1]]
}