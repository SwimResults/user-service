- different options in terms of complexity:
    - special notification settings for each combination of
        - meeting
        - athlete
        - team
    - only default settings
    - some settings for different notification types
        - files
        - favorites
        - team mates
        - myself
    - settings for time of notification → 15 min before start?

# Notify about…

- free text

## Before Meeting

- ~~new files → prob only interesting for start list? other publications?~~
- start list available/published → to find own starts and add favourites
- meeting upcoming reminder → to remember that SwimResults exists for this meeting (only when subscribed)
- new meeting with your athlete → remember user that app exists (only if not subscribed but meeting with his athlete, instead of two previous ones?)

## During Meeting

- live stream info → maybe as part of start notification (meeting has been started on time, see live timing and stream)
- me → be informed about my own start
    - upcoming start → reminder to not miss a start
    - new result published → to know places as fast as possible
    - new start imported → final qualification, changes in the start list
- favourites
    - upcoming start → don’t miss your friend’s/child’s/favourite’s starts
    - new result published → to know places as fast as possible, really necessary?
    - ~~swimming right now? → live link + stream → upcoming start is enough~~
    - new start imported → final qualification, changes in the start list
- team mates
    - ~~upcoming start → too much~~
    - ~~new result published → too much~~
    - ~~swimming right now? → live link + stream~~
    - new start imported → final qualification (only finals!)
- schedule
    - next start delayed → to know if you have more break time
    - next start early → to not miss starts
    - changes/updates (break shortened, etc.) → to know about changed breaks/warmup times

## After Meeting

- results and results file published → to receive final results without checking all the time

# Notification Settings

configure which…

Athlete Notifications

Favourites Notifications

Schedule Notifications

Meeting Notifications

… you want to receive

# Notification Service vs. User Service

| Notification Service | Notifications in User Service |
| --- | --- |
| not affected by user service availability | no new micro service |
|  | no inter-service relations since notifications are user-related |
|  | less server resources needed |
|  | user settings and notification settings in one place |
|  | closer relation to subscriptions |

# SwimResults User ↔ Notification User

- Notifications are pay feature
- with SwimResults PLUS or SwimResults PRO user has notifications
- for both an account is required
- notifications per account
- no notification user without account?
- but what about general notifications for everyone?! → notification identifier required
- what about multiple devices of the same user?
- maybe subscription not related to account?!
- app only with account?

# Questions to clarify

## How prevent multiple notifications for the same reason?

## Should subscriptions always be related to a user account?

## How to handle multiple devices of the same user?

## Are there notifications without subscription?