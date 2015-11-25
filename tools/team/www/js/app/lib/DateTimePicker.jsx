React = require('react/addons');

module.exports = React.createClass({
    getInitialState: function(){
        var now = new Date();

        var hours = now.getHours();

        var minutes = now.getMinutes();
        if (minutes >= 0 && minutes < 15) {
            minutes = 15;
        } else if (minutes >= 15 && minutes < 30) {
            minutes = 30;
        } else if (minutes >= 30 && minutes < 45) {
            minutes = 45;
        } else if (minutes >= 45 && minutes < 60) {
            minutes = "00";
            if (hours == 23) {
                hours = 0;
            } else {
                hours++;
            }
        }

        var defaultHour = hours;
        if (defaultHour == 0) {
            defaultHour = 12;
        } else if (defaultHour > 12) {
            defaultHour = defaultHour - 12;
        }

        var ampm = (now.getHours() > 11 ? "PM" : "AM");
        var timeZone = "-" + (now.getTimezoneOffset()/60) + "00";
        var dateTimeString = (now.getMonth() + 1) +  "/" + (now.getDate()) + "/" + now.getFullYear() + " " + hours + ":" + minutes  + " " + timeZone;
        var initDate = new Date(dateTimeString);

        return {
            date: initDate.getTime(),
            defaultHour: defaultHour,
            dateString: this.getDateValueString(initDate),
            hourString: defaultHour + "",
            minuteString: minutes + "",
            ampmString: ampm,
            timeZoneString: timeZone,
            timeZoneText: "",
            isDst: this.checkIfDST(now)
        }
    },
    updateDate: function(){
        var timeZoneText = this.state.timeZoneText;
        if (timeZoneText == "") {
            timeZoneText = $("#dateTimezone option:selected").text();
        }
        var tzOffset = this.state.timeZoneString;
        if (tzOffset.indexOf("-") < 0){
            tzOffset = "+" + tzOffset;
        }
        newDate = new Date(this.state.dateString + " " + this.state.hourString + ":" + this.state.minuteString + " " + this.state.ampmString + " " + tzOffset);
        this.setState(
            {
                date: newDate.getTime(),
                timeZoneText: timeZoneText,
                isDst: this.checkIfDST(newDate)
            },
            this.props.onChange);
    },
    componentDidUpdate: function(prevProps, prevState) {
        
    },
    getDateString: function(dateObj) {
        
        var days = new Array("Sun","Mon","Tue","Wed","Thu","Fri","Sat");
        var months = new Array("January", "February", "March","April", "May", "June", "July","August", "September", "October","November", "December");

        var day = days[dateObj.getDay()];
        var month = months[dateObj.getMonth()];
        var date = dateObj.getDate();
        var year = dateObj.getFullYear();

        return day + ", " + month + " " + date + ", " + year;
    },
    getDateValueString: function(dateObj) {
        return (dateObj.getMonth() +1) + "/" + dateObj.getDate() + "/" + dateObj.getFullYear();
    },
    getDates: function(){
        var dateList = new Array();
        var maxDays = this.props.maxDays;
        for (i=0;i<maxDays;i++) {
            var newDate = new Date();
            newDate.setDate(newDate.getDate() + i);
            dateList.push( <option key={i} value={this.getDateValueString(newDate)}>{this.getDateString(newDate)}</option> );
        }

        return dateList;
    },
    checkIfDST: function(date) {
        return date.getTimezoneOffset() < this.stdTimezoneOffset();
    },
    getTimeZones: function(){
        var timeZones = new Array();
        var standardUSZones = ["SAMT","MSK","EET","CET","GMT","AST","EST","CST","MST","PST","AKST","HAST"];
        var standardUSOffsets = ["400","300","200","100","000","-400","-500","-600","-700","-800","-900","-1000"];
        var daylightUSZones = ["MSD","EEST","CEST","BST","GMT","ADT","EDT","CDT","MDT","PDT","AKDT","HADT"];
        var daylightUSOffsets = ["400","300","200","100","000","-300","-400","-500","-600","-700","-800","-900"];;

        var isDst = this.state.isDst; 

        for (i=0;i<standardUSZones.length;i++){

            var value;
            var display;

            if (isDst){
                value = daylightUSOffsets[i];
                display = daylightUSZones[i];
            } else {
                value = standardUSOffsets[i];
                display = standardUSZones[i];
            }

            timeZones.push(<option key={i} value={value} >{display}</option>)
        }

        return timeZones;
    },
    stdTimezoneOffset: function() {
        var now = new Date();
        var jan = new Date(now.getFullYear(), 0, 1);
        var jul = new Date(now.getFullYear(), 6, 1);
        return Math.max(jan.getTimezoneOffset(), jul.getTimezoneOffset());
    },
    handleDateChange: function(event){
        this.setState({dateString: event.target.value},this.updateDate);
    },
    handleHourChange: function(event){
        this.setState({hourString: event.target.value},this.updateDate);
    },
    handleMinuteChange: function(event) {
        this.setState({minuteString: event.target.value},this.updateDate);
    },
    handleAmPmChange: function(event) {
        this.setState({ampmString: event.target.value},this.updateDate);
    },
    handleTimeZoneChange: function(event){
        var selectedIndex = event.target.selectedIndex;
        var timeZoneText = event.target.options[selectedIndex].text;
    
        this.setState({timeZoneString: event.target.value,timeZoneText: timeZoneText},this.updateDate);
    },
    render: function(){
        var dateList = this.getDates();
        var timeZones = this.getTimeZones();
        var now = new Date();
        var currentOffset = -(now.getTimezoneOffset()/60) + "00";

        var hourStyle = {
            paddingRight: '0px',
            // width: '6.5rem'
        }

        var minuteStyle = {
            paddingLeft: '0px',
            paddingRight: '8px',
            width: '20%'
        }

        var ampmStyle = {
            paddingLeft: '0px',
			paddingRight: '7px',
            // width: '7rem'
        }

        var timezoneStyle = {
             paddingLeft: '0px',
			 float: 'right',
        }

        var colonStyle = {
            width: '1rem',
            float: 'left',
            fontWeight: 'bold',
            fontSize: 'medium',
            padding: '2px',
            textAlign: 'center'
        }

        return(
            <div>
                <div className="form-group">
                    <label htmlFor="dateDate">Date:</label>
                    <select id="dateDate" name="dateDate" className="form-control" onChange={this.handleDateChange}>
                        {dateList}
                    </select>
                </div>
                <div className="form-group">
                    <label htmlFor="dateTime">Time:</label>
                    <div className="row">
                        <div className="col-xs-3" style={hourStyle}>
                            <select defaultValue={this.state.defaultHour} id="dateTimeHour" name="dateTimeHour" className="form-control" onChange={this.handleHourChange}>
                                <option value="1">1</option>
                                <option value="2">2</option>
                                <option value="3">3</option>
                                <option value="4">4</option>
                                <option value="5">5</option>
                                <option value="6">6</option>
                                <option value="7">7</option>
                                <option value="8">8</option>
                                <option value="9">9</option>
                                <option value="10">10</option>
                                <option value="11">11</option>
                                <option value="12">12</option>
                            </select>
                        </div>
                        <div className="col-xs-1" style={colonStyle}>:</div>
                        <div className="col-xs-3" style={minuteStyle}>
                            <select id="dateTimeMinute"  name="dateTimeMinute" className="form-control" onChange={this.handleMinuteChange} defaultValue={this.state.minuteString}>
                                <option value="00">00</option>
                                <option value="15">15</option>
                                <option value="30">30</option>
                                <option value="45">45</option>
                            </select>
                        </div>
                        <div className="col-xs-3" style={ampmStyle}>
                            <select defaultValue={this.state.ampmString} id="dateTimeAmPm"  name="dateTimeAmPm" className="form-control" onChange={this.handleAmPmChange}>
                                <option value="AM">AM</option>
                                <option value="PM">PM</option>
                            </select>
                        </div>
                        <div className="col-xs-3" style={timezoneStyle}>
                            <select id="dateTimezone" onChange={this.handleTimeZoneChange} name="dateTimeAmPm" className="form-control" defaultValue={this.state.timeZoneString}>
                                {timeZones}
                            </select>
                        </div>
                    </div>
                </div>
            </div>
            )
    }
})
