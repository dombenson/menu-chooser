module Main exposing (..)

import Html.App as App
import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (..)
import Random
import Task
import Date
import Http
import String
import Dict
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
    , eventId : Int
    }


type alias Response data =
    { errorCode : Int
    , errorMessage : String
    , data : data
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
    , date : String
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
    , options : List Int
    , optionInfo : OptionDict
    }


type alias Option =
    { id : Int
    , name : String
    , description : String
    , selected : Bool
    }


type alias FormFields =
    { loginKey : String
    }


type alias CourseDict =
    Dict.Dict Int Course


type alias OptionDict =
    Dict.Dict Int Option


type alias Model =
    { loaded : Bool
    , loggedIn : Bool
    , attendee : Attendee
    , event : Event
    , courseInfo : CourseDict
    , courses : List Int
    , form : FormFields
    }


emptyFormFields : FormFields
emptyFormFields =
    { loginKey = "" }


emptyAttendee : Attendee
emptyAttendee =
    { id = 0, name = "", email = "", eventId = 0 }


emptyEvent : Event
emptyEvent =
    { name = "", location = "", date = "" }


emptyCourse : Course
emptyCourse =
    { id = 0, name = "", options = [], optionInfo = Dict.empty }


emptyOption : Option
emptyOption =
    { id = 0, name = "", description = "", selected = False }


emptyCourses : List Int
emptyCourses =
    []


model : Model
model =
    { loaded = False
    , loggedIn = False
    , attendee = emptyAttendee
    , event = emptyEvent
    , courseInfo = (Dict.empty)
    , courses = emptyCourses
    , form = emptyFormFields
    }


init : ( Model, Cmd Msg )
init =
    ( model, loadAttendee "" )


apiEndpoint =
    "api"



-- UPDATE


type Msg
    = Noop
    | SetAttendee (Response Attendee)
    | FailAttendee Http.Error
    | FailNoop Http.Error
    | FormLoginKey String
    | DoLogin
    | SetEvent (Response Event)
    | SetCourses (Response (List BaseCourse))
    | SetOptions Int (Response (List BaseOption))
    | SetSelection Int (Response (Int))
    | TriggerSetSelection Int Int


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        Noop ->
            ( model, Cmd.none )

        FailNoop _ ->
            ( model, Cmd.none )

        SetAttendee newAtt ->
            ( { model | loaded = True, loggedIn = True, attendee = newAtt.data }, loadEvent "" )

        SetEvent newEvt ->
            ( { model | event = newEvt.data }, loadCourses "" )

        FailAttendee _ ->
            ( { model | loaded = True, loggedIn = False }, Cmd.none )

        FormLoginKey key ->
            let
                curForm =
                    model.form
            in
                ( { model | form = { curForm | loginKey = key } }, Cmd.none )

        SetCourses crsRes ->
            let
                newCourses =
                    List.map .id crsRes.data
            in
                ( { model | courses = newCourses, courseInfo = Dict.fromList (List.map makeTupleFromBase crsRes.data) }, Cmd.batch (List.map loadOptions newCourses) )

        SetOptions crsId crsRes ->
            let
                curCourses =
                    model.courseInfo
            in
                ( { model | courseInfo = Dict.update crsId (insertOpts crsRes.data) curCourses }, loadSel crsId )

        TriggerSetSelection crsId selId ->
            ( model, setSel crsId selId )

        SetSelection crsId selRes ->
            let
                curCourses =
                    model.courseInfo
            in
                ( { model | courseInfo = Dict.update crsId (insertSelection selRes.data) curCourses }, Cmd.none )

        DoLogin ->
            let
                useTok =
                    model.form.loginKey
            in
                ( { model | loaded = False }, loginAttendee useTok )


runUpdate : msg -> Cmd msg
runUpdate toExec =
    Task.perform identity identity (Task.succeed toExec)


makeTupleFromBase : BaseCourse -> ( Int, Course )
makeTupleFromBase bse =
    ( bse.id, { emptyCourse | id = bse.id, name = bse.name } )


makeOptionTupleFromBase : BaseOption -> ( Int, Option )
makeOptionTupleFromBase bse =
    ( bse.id, { id = bse.id, name = bse.name, description = bse.description, selected = False } )


makeOptFromBase : BaseOption -> Option
makeOptFromBase bse =
    { id = bse.id, name = bse.name, description = bse.description, selected = False }


insertOpts : List BaseOption -> Maybe Course -> Maybe Course
insertOpts opts mbeCourse =
    case mbeCourse of
        Just course ->
            Just { course | options = List.map .id opts, optionInfo = Dict.fromList (List.map makeOptionTupleFromBase opts) }

        Nothing ->
            Nothing


selectIfMatch : Int -> Int -> Option -> Option
selectIfMatch selId curId curOpt =
    { curOpt | selected = (selId == curId) }


insertSelection : Int -> Maybe Course -> Maybe Course
insertSelection selId mbeCourse =
    case mbeCourse of
        Just course ->
            let
                curOpts =
                    course.optionInfo
            in
                Just { course | optionInfo = Dict.map (selectIfMatch selId) curOpts }

        Nothing ->
            Nothing



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
attendeeResponseDecoder =
    responseDecoder attendeeDecoder


eventResponseDecoder : Json.Decoder (Response Event)
eventResponseDecoder =
    responseDecoder eventDecoder


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


selectionResponseDecoder : Json.Decoder (Response (Int))
selectionResponseDecoder =
    Json.object3 Response
        ("errorCode" := Json.int)
        ("errorMessage" := string)
        ("data" := Json.int)


subscriptions : Model -> Sub Msg
subscriptions model =
    Sub.none


loginAttendee : String -> Cmd Msg
loginAttendee tok =
    let
        url =
            apiEndpoint ++ "/login"
    in
        Task.perform FailAttendee SetAttendee (Http.post attendeeResponseDecoder url (Http.string tok))


loadAttendee : String -> Cmd Msg
loadAttendee msg =
    let
        url =
            apiEndpoint ++ "/user/me"
    in
        Task.perform FailAttendee SetAttendee (Http.get attendeeResponseDecoder url)


loadEvent : String -> Cmd Msg
loadEvent msg =
    let
        url =
            apiEndpoint ++ "/user/event"
    in
        Task.perform FailAttendee SetEvent (Http.get eventResponseDecoder url)


loadCourses : String -> Cmd Msg
loadCourses msg =
    let
        url =
            apiEndpoint ++ "/user/courses"
    in
        Task.perform FailAttendee SetCourses (Http.get coursesResponseDecoder url)


loadOptions : Int -> Cmd Msg
loadOptions crsId =
    let
        url =
            apiEndpoint ++ "/user/options/" ++ (toString crsId)
    in
        Task.perform FailAttendee (SetOptions crsId) (Http.get optionsResponseDecoder url)


loadSel : Int -> Cmd Msg
loadSel crsId =
    let
        url =
            apiEndpoint ++ "/user/selection/" ++ (toString crsId)
    in
        Task.perform FailNoop (SetSelection crsId) (Http.get selectionResponseDecoder url)


setSel : Int -> Int -> Cmd Msg
setSel crsId selId =
    let
        url =
            apiEndpoint ++ "/user/selection/" ++ (toString crsId)
    in
        Task.perform FailNoop (SetSelection crsId) (Http.post selectionResponseDecoder url (Http.string (toString selId)))



-- VIEW


drawOption : Int -> OptionDict -> Int -> Html Msg
drawOption crsId opts optId =
    let
        opt =
            Maybe.withDefault emptyOption (Dict.get optId opts)

        selectThis =
            TriggerSetSelection crsId optId
    in
        let
            className =
                if opt.selected then
                    "selected option"
                else
                    "option"
        in
            div [ onClick selectThis, class className ]
                [ div [ class "name" ] [ text opt.name ]
                , div [ class "description" ] [ text opt.description ]
                , div [ class "saved" ] [ text "Saved" ]
                ]


drawCourse : CourseDict -> Int -> Html Msg
drawCourse crses crsId =
    let
        crs =
            Maybe.withDefault emptyCourse (Dict.get crsId crses)

        optRendr =
            drawOption crsId crs.optionInfo
    in
        div [ class "course" ]
            [ div [ class "name" ] [ text crs.name ]
            , div [ class "options" ] (List.map optRendr crs.options)
            ]


drawNotLoggedIn : Model -> Html Msg
drawNotLoggedIn model =
    div [ class "login" ]
        [ div [ class "title" ] [ text "Please log in" ]
        , div [ class "body" ] [ div [ class "info" ] [ text "Enter the code from your invitation email" ]
        , div []
            [ input [ type' "text", onInput FormLoginKey ] []
            ]
        , div []
            [ button [ onClick DoLogin ] [ text "Login" ]
            ]
        ] ]


drawLoading : Model -> Html a
drawLoading model =
    div []
        [ div [] [ text "Loading" ]
        ]


view : Model -> Html Msg
view model =
    if model.loaded then
        if model.loggedIn then
            let
                crsRendr =
                    drawCourse model.courseInfo

                evtDate =
                    Date.fromString model.event.date |> Result.withDefault (Date.fromTime 0)
            in
                div [ class "loaded" ]
                    [ div [ id "header" ]
                        [ div [ class "attendee" ] [ span [ class "name" ] [ text model.attendee.name ] ]
                        , div [ class "event" ] [ span [ class "name" ] [ text model.event.name ], span [ class "location" ] [ text model.event.location ], span [ class "date" ] [ text ((toString (Date.day evtDate)) ++ " " ++ (toString (Date.month evtDate)) ++ " " ++ (toString (Date.year evtDate))) ], span [ class "time" ] [ text ((toString (Date.hour evtDate)) ++ ":" ++ (String.pad 2 '0' (toString (Date.minute evtDate)))) ] ]
                        ]
                    , div [ id "body" ]
                        [ div [ id "courses" ] (List.map crsRendr model.courses)
                        ]
                    ]
        else
            drawNotLoggedIn model
    else
        drawLoading model
