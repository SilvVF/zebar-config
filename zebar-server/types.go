package main

import "time"

type DailyNoteCommon struct {
	Game             GameId
	Current          int
	Max              int
	FullyRecoveredTs int
	RecoverInterval  time.Duration
}

type DailyNoteResponseStarRail struct {
	Retcode int    `json:"retcode"`
	Message string `json:"message"`
	Data    struct {
		CurrentStamina       int `json:"current_stamina"`
		MaxStamina           int `json:"max_stamina"`
		StaminaRecoverTime   int `json:"stamina_recover_time"`
		StaminaFullTs        int `json:"stamina_full_ts"`
		AcceptedEpeditionNum int `json:"accepted_epedition_num"`
		TotalExpeditionNum   int `json:"total_expedition_num"`
		Expeditions          []struct {
			Avatars       []string `json:"avatars"`
			Status        string   `json:"status"`
			RemainingTime int      `json:"remaining_time"`
			Name          string   `json:"name"`
			ItemURL       string   `json:"item_url"`
			FinishTs      int      `json:"finish_ts"`
		} `json:"expeditions"`
		CurrentTrainScore        int  `json:"current_train_score"`
		MaxTrainScore            int  `json:"max_train_score"`
		CurrentRogueScore        int  `json:"current_rogue_score"`
		MaxRogueScore            int  `json:"max_rogue_score"`
		WeeklyCocoonCnt          int  `json:"weekly_cocoon_cnt"`
		WeeklyCocoonLimit        int  `json:"weekly_cocoon_limit"`
		CurrentReserveStamina    int  `json:"current_reserve_stamina"`
		IsReserveStaminaFull     bool `json:"is_reserve_stamina_full"`
		RogueTournWeeklyUnlocked bool `json:"rogue_tourn_weekly_unlocked"`
		RogueTournWeeklyMax      int  `json:"rogue_tourn_weekly_max"`
		RogueTournWeeklyCur      int  `json:"rogue_tourn_weekly_cur"`
		CurrentTs                int  `json:"current_ts"`
		RogueTournExpIsFull      bool `json:"rogue_tourn_exp_is_full"`
	} `json:"data"`
}

type DailyNoteResponseZZZ struct {
	Retcode int    `json:"retcode"`
	Message string `json:"message"`
	Data    struct {
		Energy struct {
			Progress struct {
				Max     int `json:"max"`
				Current int `json:"current"`
			} `json:"progress"`
			Restore int `json:"restore"`
			DayType int `json:"day_type"`
			Hour    int `json:"hour"`
			Minute  int `json:"minute"`
		} `json:"energy"`
		Vitality struct {
			Max     int `json:"max"`
			Current int `json:"current"`
		} `json:"vitality"`
		VhsSale struct {
			SaleState string `json:"sale_state"`
		} `json:"vhs_sale"`
		CardSign         string `json:"card_sign"`
		BountyCommission struct {
			Num         int `json:"num"`
			Total       int `json:"total"`
			RefreshTime int `json:"refresh_time"`
		} `json:"bounty_commission"`
		SurveyPoints any `json:"survey_points"`
		AbyssRefresh int `json:"abyss_refresh"`
		Coffee       any `json:"coffee"`
		WeeklyTask   struct {
			RefreshTime int `json:"refresh_time"`
			CurPoint    int `json:"cur_point"`
			MaxPoint    int `json:"max_point"`
		} `json:"weekly_task"`
		MemberCard struct {
			IsOpen          bool   `json:"is_open"`
			MemberCardState string `json:"member_card_state"`
			ExpTime         string `json:"exp_time"`
		} `json:"member_card"`
		IsSub      bool `json:"is_sub"`
		IsOtherSub bool `json:"is_other_sub"`
	} `json:"data"`
}

type DailyNoteResponseGenshin struct {
	Retcode int    `json:"retcode"`
	Message string `json:"message"`
	Data    struct {
		CurrentResin              int    `json:"current_resin"`
		MaxResin                  int    `json:"max_resin"`
		ResinRecoveryTime         string `json:"resin_recovery_time"`
		FinishedTaskNum           int    `json:"finished_task_num"`
		TotalTaskNum              int    `json:"total_task_num"`
		IsExtraTaskRewardReceived bool   `json:"is_extra_task_reward_received"`
		RemainResinDiscountNum    int    `json:"remain_resin_discount_num"`
		ResinDiscountNumLimit     int    `json:"resin_discount_num_limit"`
		CurrentExpeditionNum      int    `json:"current_expedition_num"`
		MaxExpeditionNum          int    `json:"max_expedition_num"`
		Expeditions               []struct {
			AvatarSideIcon string `json:"avatar_side_icon"`
			Status         string `json:"status"`
			RemainedTime   string `json:"remained_time"`
		} `json:"expeditions"`
		CurrentHomeCoin      int    `json:"current_home_coin"`
		MaxHomeCoin          int    `json:"max_home_coin"`
		HomeCoinRecoveryTime string `json:"home_coin_recovery_time"`
		CalendarURL          string `json:"calendar_url"`
		Transformer          struct {
			Obtained     bool `json:"obtained"`
			RecoveryTime struct {
				Day     int  `json:"Day"`
				Hour    int  `json:"Hour"`
				Minute  int  `json:"Minute"`
				Second  int  `json:"Second"`
				Reached bool `json:"reached"`
			} `json:"recovery_time"`
			Wiki        string `json:"wiki"`
			Noticed     bool   `json:"noticed"`
			LatestJobID string `json:"latest_job_id"`
		} `json:"transformer"`
		DailyTask struct {
			TotalNum                  int  `json:"total_num"`
			FinishedNum               int  `json:"finished_num"`
			IsExtraTaskRewardReceived bool `json:"is_extra_task_reward_received"`
			TaskRewards               []struct {
				Status string `json:"status"`
			} `json:"task_rewards"`
			AttendanceRewards []struct {
				Status   string `json:"status"`
				Progress int    `json:"progress"`
			} `json:"attendance_rewards"`
			AttendanceVisible                bool   `json:"attendance_visible"`
			StoredAttendance                 string `json:"stored_attendance"`
			StoredAttendanceRefreshCountdown int    `json:"stored_attendance_refresh_countdown"`
		} `json:"daily_task"`
		ArchonQuestProgress struct {
			List                    []interface{} `json:"list"`
			IsOpenArchonQuest       bool          `json:"is_open_archon_quest"`
			IsFinishAllMainline     bool          `json:"is_finish_all_mainline"`
			IsFinishAllInterchapter bool          `json:"is_finish_all_interchapter"`
			WikiURL                 string        `json:"wiki_url"`
		} `json:"archon_quest_progress"`
	} `json:"data"`
}
