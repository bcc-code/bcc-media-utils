# Set project and dataset
PROJECT_ID="bccm-k8s-main"
DATASET="rudderstack_prod"

# List of views
declare -a VIEWS=(
    "achievement_clicked_view"
    "achievement_shared_view"
    "application_backgrounded_view"
    "application_installed_all_view"
    "application_installed_view"
    "application_open_view_model"
    "application_opened_view"
    "application_updated_view"
    "audioonly_clicked_view"
    "calendarday_clicked_view"
    "chapter_clicked_view"
    "content_shared_view"
    "deep_link_opened_view"
    "download_started_view"
    "episode_download_view"
    "episodes_view"
    "game_closed_all_view"
    "game_closed_view"
    "identifies_view"
    "interaction_view"
    "language_changed_view"
    "pages_view"
    "rudder_discards_view"
    "screens_view"
    "search_performed_view"
)

# Loop through each view
for VIEW_NAME in "${VIEWS[@]}"
do
    # Get view definition
    VIEW_QUERY=$(bq show --format=prettyjson $PROJECT_ID:$DATASET.$VIEW_NAME | jq -r '.view.query')

    # Delete the view
    bq rm -f $PROJECT_ID:$DATASET.$VIEW_NAME

    # Recreate the view
    bq mk --use_legacy_sql=false --view "$VIEW_QUERY" $PROJECT_ID:$DATASET.$VIEW_NAME
done

