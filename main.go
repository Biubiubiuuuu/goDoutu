package main

const (
    Url = "https://www.zhihu.com/api/v4/questions/%v/answers"
    Female = 1
    Male = 2
)

// 知乎问题回答返回结构体
type ResponseData struct {
	Data []struct {
		ID int `json:"id"`
		Type string `json:"type"`
		AnswerType string `json:"answer_type"`
		Question struct {
			Type string `json:"type"`
			ID int `json:"id"`
			Title string `json:"title"`
			QuestionType string `json:"question_type"`
			Created int `json:"created"`
			UpdatedTime int `json:"updated_time"`
			URL string `json:"url"`
			Relationship struct {
			} `json:"relationship"`
		} `json:"question"`
		Author struct {
			ID string `json:"id"`
			URLToken string `json:"url_token"`
			Name string `json:"name"`
			AvatarURL string `json:"avatar_url"`
			AvatarURLTemplate string `json:"avatar_url_template"`
			IsOrg bool `json:"is_org"`
			Type string `json:"type"`
			URL string `json:"url"`
			UserType string `json:"user_type"`
			Headline string `json:"headline"`
			Badge []interface{} `json:"badge"`
			BadgeV2 struct {
				Title string `json:"title"`
				MergedBadges []interface{} `json:"merged_badges"`
				DetailBadges []interface{} `json:"detail_badges"`
				Icon string `json:"icon"`
				NightIcon string `json:"night_icon"`
			} `json:"badge_v2"`
			Gender int `json:"gender"`
			IsAdvertiser bool `json:"is_advertiser"`
			FollowerCount int `json:"follower_count"`
			IsFollowed bool `json:"is_followed"`
			IsPrivacy bool `json:"is_privacy"`
		} `json:"author"`
		URL string `json:"url"`
		IsCollapsed bool `json:"is_collapsed"`
		CreatedTime int `json:"created_time"`
		UpdatedTime int `json:"updated_time"`
		Extras string `json:"extras"`
		IsCopyable bool `json:"is_copyable"`
		IsNormal bool `json:"is_normal"`
		VoteupCount int `json:"voteup_count"`
		CommentCount int `json:"comment_count"`
		IsSticky bool `json:"is_sticky"`
		AdminClosedComment bool `json:"admin_closed_comment"`
		CommentPermission string `json:"comment_permission"`
		CanComment struct {
			Reason string `json:"reason"`
			Status bool `json:"status"`
		} `json:"can_comment"`
		ReshipmentSettings string `json:"reshipment_settings"`
		Content string `json:"content"`
		EditableContent string `json:"editable_content"`
		Excerpt string `json:"excerpt"`
		CollapsedBy string `json:"collapsed_by"`
		CollapseReason string `json:"collapse_reason"`
		AnnotationAction interface{} `json:"annotation_action"`
		MarkInfos []interface{} `json:"mark_infos"`
		RelevantInfo struct {
			IsRelevant bool `json:"is_relevant"`
			RelevantType string `json:"relevant_type"`
			RelevantText string `json:"relevant_text"`
		} `json:"relevant_info"`
		SuggestEdit struct {
			Reason string `json:"reason"`
			Status bool `json:"status"`
			Tip string `json:"tip"`
			Title string `json:"title"`
			UnnormalDetails struct {
				Status string `json:"status"`
				Description string `json:"description"`
				Reason string `json:"reason"`
				ReasonID int `json:"reason_id"`
				Note string `json:"note"`
			} `json:"unnormal_details"`
			URL string `json:"url"`
		} `json:"suggest_edit"`
		IsLabeled bool `json:"is_labeled"`
		RewardInfo struct {
			CanOpenReward bool `json:"can_open_reward"`
			IsRewardable bool `json:"is_rewardable"`
			RewardMemberCount int `json:"reward_member_count"`
			RewardTotalMoney int `json:"reward_total_money"`
			Tagline string `json:"tagline"`
		} `json:"reward_info"`
		Relationship struct {
			IsAuthor bool `json:"is_author"`
			IsAuthorized bool `json:"is_authorized"`
			IsNothelp bool `json:"is_nothelp"`
			IsThanked bool `json:"is_thanked"`
			IsRecognized bool `json:"is_recognized"`
			Voting int `json:"voting"`
			UpvotedFollowees []interface{} `json:"upvoted_followees"`
		} `json:"relationship"`
		AdAnswer interface{} `json:"ad_answer"`
	} `json:"data"`
	Paging struct {
		IsEnd bool `json:"is_end"`
		IsStart bool `json:"is_start"`
		Next string `json:"next"`
		Previous string `json:"previous"`
		Totals int `json:"totals"`
	} `json:"paging"`
}

func main() {
  
}