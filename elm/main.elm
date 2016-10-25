import Html.App as App
import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (..)
import Random
import Task
import Date
import Http
import Json.Decode as Json exposing (Decoder, decodeValue, succeed, string, oneOf, null, list, bool, (:=), andThen)

main =
    App.program
        { init = init
        , view = view
        , update = update
        , subscriptions = subscriptions
        }



-- MODEL


type alias Attendee =
    { id : Int
    , name : String
    , email : String
    , eventId: Int
    }

type alias Response data =
    { errorCode: Int
    , errorMessage: String
    , data: data
    }

{-
-- an alternate way of writing Response to capture whether an error occurred or not
-- in the types:
type alias AltResponse data =
    { error : ResponseError
    , data : data
    }

-- used in our alt response; the error is either NoError or an Error with
-- an int code and a string message:
type ResponseError = NoError | Error Int String
-}

type alias Event =
    { name : String
    , location : String
    , time : String
    }

type alias BaseCourse =
    { id : Int
    , name : String
    }

type alias BaseOption =
    { id : Int
    , name : String
    , description : String
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
    { id = 0, name = "", email = "" , eventId = 0}


emptyEvent : Event
emptyEvent =
    { name = "", location = "", time = "" }


emptyCourses : List Course
emptyCourses =
    []


model : Model
model =
    { loaded = False
    , loggedIn = False
    , attendee = emptyAttendee
    , event = emptyEvent
    , courses = emptyCourses
    , form = emptyFormFields }


init : ( Model, Cmd Msg )
init =
    ( model, loadAttendee "" )

apiEndpoint = "/api"



-- UPDATE


type Msg
    = Noop
    | SetAttendee (Response Attendee)
    | FailAttendee Http.Error
    | FormLoginKey String
    | DoLogin
    | SetEvent (Response Event)
    | SetCourses (Response (List BaseCourse))
    | SetOptions Int (Response (List BaseOption))


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        Noop ->
            ( model, Cmd.none )
        SetAttendee newAtt ->
            ( {model | loaded = True, loggedIn = True, attendee = newAtt.data}, loadEvent "")
        SetEvent newEvt ->
            ( {model | event = newEvt.data}, loadCourses "")
        FailAttendee _ ->
            ( {model | loaded = True, loggedIn = False}, Cmd.none)
        FormLoginKey key ->
            let curForm = model.form
            in ( {model | form = {curForm | loginKey = key}}, Cmd.none)
        SetCourses crsRes ->
            ( {model | courses = List.map makeFromBase crsRes.data}, Cmd.none)
        SetOptions crsId crsRes ->
            ( model, Cmd.none)
        DoLogin ->
            let useTok = model.form.loginKey
            in ( {model | loaded = False}, loginAttendee useTok )


runUpdate : msg -> Cmd msg
runUpdate toExec =
    Task.perform identity identity (Task.succeed toExec)

makeFromBase : BaseCourse -> Course
makeFromBase bse = { id = bse.id, name = bse.name, options = []}



-- SUBSCRIPTIONS

responseDecoder : Json.Decoder data -> Json.Decoder (Response data)
responseDecoder dataDecoder =
    Json.object3 Response
        ("errorCode" := Json.int)
        ("errorMessage" := Json.string)
        ("data" := dataDecoder)

attendeeDecoder : Json.Decoder Attendee
attendeeDecoder =
    Json.object4 Attendee
        ("id" := Json.int)
        ("name" := string)
        ("email" := string)
        ("eventId" := Json.int)

eventDecoder : Json.Decoder Event
eventDecoder =
    Json.object3 Event
        ("name" := string)
        ("location" := string)
        ("date" := string)

attendeeResponseDecoder : Json.Decoder (Response Attendee)
attendeeResponseDecoder = responseDecoder attendeeDecoder

eventResponseDecoder : Json.Decoder (Response Event)
eventResponseDecoder = responseDecoder eventDecoder

baseCourseDecoder : Json.Decoder BaseCourse
baseCourseDecoder =
    Json.object2 BaseCourse
        ("id" := Json.int)
        ("name" := string)

coursesResponseDecoder : Json.Decoder (Response (List BaseCourse))
coursesResponseDecoder =
    Json.object3 Response
        ("errorCode" := Json.int)
        ("errorMessage" := string)
        ("data" := (Json.list baseCourseDecoder))


baseOptionDecoder : Json.Decoder BaseOption
baseOptionDecoder =
    Json.object3 BaseOption
        ("id" := Json.int)
        ("name" := string)
        ("description" := string)

optionsResponseDecoder : Json.Decoder (Response (List BaseOption))
optionsResponseDecoder =
    Json.object3 Response
        ("errorCode" := Json.int)
        ("errorMessage" := string)
        ("data" := (Json.list baseOptionDecoder))


subscriptions : Model -> Sub Msg
subscriptions model =
    Sub.none

loginAttendee : String -> Cmd Msg
loginAttendee tok =
    let url = apiEndpoint ++ "/login/" ++ tok
    in
        Task.perform FailAttendee SetAttendee (Http.get attendeeResponseDecoder url)

loadAttendee : String -> Cmd Msg
loadAttendee msg =
    let url = apiEndpoint ++ "/user/me"
    in
        Task.perform FailAttendee SetAttendee (Http.get attendeeResponseDecoder url)

loadEvent : String -> Cmd Msg
loadEvent msg =
    let url = apiEndpoint ++ "/user/event"
    in
        Task.perform FailAttendee SetEvent (Http.get eventResponseDecoder url)

loadCourses : String -> Cmd Msg
loadCourses msg =
    let url = apiEndpoint ++ "/user/courses"
    in
        Task.perform FailAttendee SetCourses (Http.get coursesResponseDecoder url)

loadOptions : Course -> Cmd Msg
loadOptions crs =
    let url = apiEndpoint ++ "/user/options/" ++ (toString crs.id)
    in
        Task.perform FailAttendee (SetOptions crs.id) (Http.get optionsResponseDecoder url)

-- VIEW

drawOption : Option -> Html a
drawOption opt =
    div [] [
        div [] [ text opt.name ],
        div [] [ text opt.description ]
    ]


drawCourse : Course -> Html a
drawCourse crs =
    div [] [
        div [] [ text crs.name ],
        div [] (List.map drawOption crs.options)
    ]

view : Model -> Html Msg
view model =
    if model.loaded then
        if model.loggedIn then
            div []
                [
                 div [] [ text model.attendee.name, text model.attendee.email ],
                 div [] [ text model.event.name, text model.event.location ],
                 div [] (List.map drawCourse model.courses)
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
