import Html.App as App
import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (..)
import Random
import Task
import Date
import Http
import Json.Decode as Json


main =
    App.program
        { init = init
        , view = view
        , update = update
        , subscriptions = subscriptions
        }



-- MODEL


type alias Attendee =
    { name : String
    , email : String
    }


type alias Event =
    { name : String
    , location : String
    , time : String
    }


type alias Course =
    { id : Int
    , name : String
    , options : List Option
    }


type alias Option =
    { id : Int
    , name : String
    , description : String
    , selected : Bool
    }

type alias FormFields = {
    loginKey: String
    }

type alias Model =
    { loaded : Bool
    , loggedIn : Bool
    , attendee : Attendee
    , event : Event
    , courses : List Course
    , form: FormFields
    }

emptyFormFields : FormFields
emptyFormFields = {loginKey = ""}


emptyAttendee : Attendee
emptyAttendee =
    { name = "", email = "" }


emptyEvent : Event
emptyEvent =
    { name = "", location = "", time = "" }


emptyCourses : List Course
emptyCourses =
    []


model : Model
model =
    { loaded = False, loggedIn = False, attendee = emptyAttendee, event = emptyEvent, courses = emptyCourses, form = emptyFormFields }


init : ( Model, Cmd Msg )
init =
    ( model, loadAttendee "" )

apiEndpoint = "http://localhost:8000"



-- UPDATE


type Msg
    = Noop | SetAttendee String | FailAttendee Http.Error | FormLoginKey String | DoLogin


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        Noop ->
            ( model, Cmd.none )
        SetAttendee newAttName ->
            ( {model | loaded = True, loggedIn = True, attendee = {emptyAttendee | name = newAttName}}, Cmd.none)
        FailAttendee _ ->
            ( {model | loaded = True, loggedIn = False}, Cmd.none)
        FormLoginKey key ->
            let curForm = model.form
            in ( {model | form = {curForm | loginKey = key}}, Cmd.none)
        DoLogin ->
            let useTok = model.form.loginKey
            in ( {model | loaded = False}, loginAttendee useTok )


runUpdate : msg -> Cmd msg
runUpdate toExec =
    Task.perform identity identity (Task.succeed toExec)



-- SUBSCRIPTIONS


subscriptions : Model -> Sub Msg
subscriptions model =
    Sub.none

loginAttendee : String -> Cmd Msg
loginAttendee tok =
    let url = apiEndpoint ++ "/login/" ++ tok
    in
        Task.perform FailAttendee SetAttendee (Http.get decodeAttendee url)

loadAttendee : String -> Cmd Msg
loadAttendee msg =
    let url = apiEndpoint ++ "/user/me"
    in
        Task.perform FailAttendee SetAttendee (Http.get decodeAttendee url)

decodeAttendee : Json.Decoder String
decodeAttendee =
    Json.at ["data", "name"] Json.string

-- VIEW


view : Model -> Html Msg
view model =
    if model.loaded then
        if model.loggedIn then
            div []
                [ div [] [ text (toString model.attendee.name), text (toString model.attendee.email) ]
                ]
        else
            div [] [
                div[] [ text "Please log in" ],
                div[] [
                    input [ type' "text", onInput FormLoginKey ] []
                    ],
                div[] [
                    button [onClick DoLogin ] [ text "Login" ]
                    ]
            ]
    else
        div []
            [ div [] [ text "Loading" ]
            ]
